package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// ScoreService computes and serves builder reputation scores.
type ScoreService struct {
	scores   domain.ScoreRepository
	projects domain.ProjectRepository
	users    domain.UserRepository
	stripe   domain.StripeRepository
	enq      ScoreEnqueuer
}

// NewScoreService constructs a ScoreService.
func NewScoreService(
	scores domain.ScoreRepository,
	projects domain.ProjectRepository,
	users domain.UserRepository,
	stripe domain.StripeRepository,
	enq ScoreEnqueuer,
) *ScoreService {
	return &ScoreService{scores: scores, projects: projects, users: users, stripe: stripe, enq: enq}
}

// GetByUsername returns the stored score for a builder, computing it on the fly
// if none exists yet.
func (s *ScoreService) GetByUsername(ctx context.Context, username string) (*domain.BuilderScore, error) {
	u, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("ScoreService.GetByUsername: %w", err)
	}
	score, err := s.scores.GetByBuilderID(ctx, u.ID)
	if err == nil {
		return score, nil
	}
	if !isNotFound(err) {
		return nil, fmt.Errorf("ScoreService.GetByUsername: %w", err)
	}
	// No score yet — compute and persist one now.
	return s.Recalculate(ctx, u.ID)
}

// Recalculate recomputes a builder's score from current project, review, and
// Stripe data, persists it, and returns it. This is invoked both synchronously
// and from the async worker.
func (s *ScoreService) Recalculate(ctx context.Context, builderID string) (*domain.BuilderScore, error) {
	projects, err := s.projects.ListByBuilder(ctx, builderID)
	if err != nil {
		return nil, fmt.Errorf("ScoreService.Recalculate: %w", err)
	}

	inputs := s.aggregate(ctx, builderID, projects)

	score := domain.ComputeScore(inputs)
	score.BuilderID = builderID
	score.CalculatedAt = time.Now().UTC()

	if err := s.scores.Upsert(ctx, &score); err != nil {
		return nil, fmt.Errorf("ScoreService.Recalculate: %w", err)
	}
	return &score, nil
}

// aggregate rolls a builder's projects and Stripe snapshot into ScoreInputs.
func (s *ScoreService) aggregate(ctx context.Context, builderID string, projects []*domain.Project) domain.ScoreInputs {
	var (
		verifiedUsers int
		weightedRating float64
		reviewCount   int
	)
	for _, p := range projects {
		if p.Status == domain.ProjectStatusArchived {
			continue
		}
		verifiedUsers += p.VerifiedUsers
		weightedRating += p.AverageRating * float64(p.ReviewCount)
		reviewCount += p.ReviewCount
	}

	avgRating := 0.0
	if reviewCount > 0 {
		avgRating = weightedRating / float64(reviewCount)
	}

	in := domain.ScoreInputs{
		VerifiedUsers: verifiedUsers,
		AverageRating: avgRating,
		ReviewCount:   reviewCount,
	}

	// Revenue and Stripe verification come from the latest Stripe snapshot.
	snap, err := s.stripe.Latest(ctx, builderID)
	if err == nil && snap != nil {
		in.MRR = snap.MRR
		in.StripeVerified = true
		// Velocity: treat current customers as the growth proxy until we have
		// a trailing 90-day comparison snapshot.
		in.NinetyDayGrowth = float64(snap.TotalCustomers)
	} else if err != nil && !errors.Is(err, domain.ErrNotFound) {
		// Non-fatal: fall back to unverified scoring.
		_ = err
	}
	return in
}

// EnqueueRecalc schedules an async recalculation for a builder.
func (s *ScoreService) EnqueueRecalc(ctx context.Context, builderID string) error {
	if s.enq == nil {
		return fmt.Errorf("ScoreService.EnqueueRecalc: %w", domain.ErrInternal)
	}
	if err := s.enq.EnqueueScoreRecalc(ctx, builderID); err != nil {
		return fmt.Errorf("ScoreService.EnqueueRecalc: %w", err)
	}
	return nil
}

// Leaderboard returns the top builders by total score.
func (s *ScoreService) Leaderboard(ctx context.Context, limit int) ([]*domain.BuilderScore, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	out, err := s.scores.TopBuilders(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("ScoreService.Leaderboard: %w", err)
	}
	return out, nil
}
