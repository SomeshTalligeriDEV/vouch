package service

import (
	"context"
	"fmt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// ProblemService holds business logic for the demand board.
type ProblemService struct {
	problems domain.ProblemRepository
	emailEnq EmailEnqueuer
}

// NewProblemService constructs a ProblemService.
func NewProblemService(problems domain.ProblemRepository, emailEnq EmailEnqueuer) *ProblemService {
	return &ProblemService{problems: problems, emailEnq: emailEnq}
}

// ProblemInput carries the fields needed to post a problem.
type ProblemInput struct {
	Title       string
	Description string
	Tags        []string
	BudgetMin   float64
	BudgetMax   float64
}

// Create posts a new problem to the demand board.
func (s *ProblemService) Create(ctx context.Context, posterID string, in ProblemInput) (*domain.Problem, error) {
	if in.BudgetMax > 0 && in.BudgetMax < in.BudgetMin {
		return nil, fmt.Errorf("ProblemService.Create: %w", domain.ErrInvalidInput)
	}
	p := &domain.Problem{
		PosterID:    posterID,
		Title:       in.Title,
		Slug:        slugify(in.Title),
		Description: in.Description,
		Tags:        in.Tags,
		BudgetMin:   in.BudgetMin,
		BudgetMax:   in.BudgetMax,
		Status:      domain.ProblemStatusOpen,
		UpvotedBy:   []string{},
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}
	if err := s.problems.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("ProblemService.Create: %w", err)
	}
	return p, nil
}

// Get returns a single problem by id.
func (s *ProblemService) Get(ctx context.Context, id string) (*domain.Problem, error) {
	p, err := s.problems.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProblemService.Get: %w", err)
	}
	return p, nil
}

// List returns a filtered, paginated set of problems.
func (s *ProblemService) List(ctx context.Context, f domain.ProblemFilter) ([]*domain.Problem, int64, error) {
	normalizePage(&f.Page, &f.Limit)
	out, total, err := s.problems.List(ctx, f)
	if err != nil {
		return nil, 0, fmt.Errorf("ProblemService.List: %w", err)
	}
	return out, total, nil
}

// Claim lets a builder claim an open problem. A builder cannot claim their own
// problem, and only one builder can win the claim (enforced atomically).
func (s *ProblemService) Claim(ctx context.Context, id, builderID string) (*domain.Problem, error) {
	p, err := s.problems.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProblemService.Claim: %w", err)
	}
	if p.PosterID == builderID {
		return nil, fmt.Errorf("ProblemService.Claim: %w", domain.ErrSelfClaim)
	}
	claimed, err := s.problems.Claim(ctx, id, builderID)
	if err != nil {
		return nil, fmt.Errorf("ProblemService.Claim: %w", err)
	}
	// Best-effort: notify the poster asynchronously.
	if s.emailEnq != nil {
		_ = s.emailEnq.EnqueueProblemClaimedEmail(ctx, claimed.ID)
	}
	return claimed, nil
}

// Upvote records a unique upvote from a user.
func (s *ProblemService) Upvote(ctx context.Context, id, userID string) (*domain.Problem, error) {
	p, err := s.problems.AddUpvote(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("ProblemService.Upvote: %w", err)
	}
	return p, nil
}
