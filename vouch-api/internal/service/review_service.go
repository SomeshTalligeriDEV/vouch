package service

import (
	"context"
	"fmt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// ReviewService holds business logic for reviews. Posting a review refreshes
// the project's rating stats and schedules a score recalculation for the
// owning builder.
type ReviewService struct {
	reviews  domain.ReviewRepository
	projects domain.ProjectRepository
	users    domain.UserRepository
	scoreEnq ScoreEnqueuer
}

// NewReviewService constructs a ReviewService.
func NewReviewService(
	reviews domain.ReviewRepository,
	projects domain.ProjectRepository,
	users domain.UserRepository,
	scoreEnq ScoreEnqueuer,
) *ReviewService {
	return &ReviewService{reviews: reviews, projects: projects, users: users, scoreEnq: scoreEnq}
}

// ReviewInput carries the fields needed to post a review.
type ReviewInput struct {
	ProjectID string
	Rating    int
	Body      string
}

// Create posts a review. A reviewer may review a project only once, and may not
// review their own project.
func (s *ReviewService) Create(ctx context.Context, reviewerID string, in ReviewInput) (*domain.Review, error) {
	if in.Rating < 1 || in.Rating > 5 {
		return nil, fmt.Errorf("ReviewService.Create: %w", domain.ErrInvalidInput)
	}
	project, err := s.projects.GetByID(ctx, in.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("ReviewService.Create: %w", err)
	}
	if project.BuilderID == reviewerID {
		return nil, fmt.Errorf("ReviewService.Create: %w", domain.ErrForbidden)
	}
	if _, err := s.reviews.GetByProjectAndReviewer(ctx, in.ProjectID, reviewerID); err == nil {
		return nil, fmt.Errorf("ReviewService.Create: %w", domain.ErrAlreadyExists)
	} else if !isNotFound(err) {
		return nil, fmt.Errorf("ReviewService.Create: %w", err)
	}

	reviewer, err := s.users.GetByID(ctx, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("ReviewService.Create: %w", err)
	}

	r := &domain.Review{
		ProjectID:        in.ProjectID,
		ReviewerID:       reviewerID,
		ReviewerUsername: reviewer.Username,
		Rating:           in.Rating,
		Body:             in.Body,
		VerifiedPurchase: false,
	}
	if err := s.reviews.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("ReviewService.Create: %w", err)
	}

	// Refresh denormalized rating stats on the project.
	stats, err := s.reviews.StatsForProject(ctx, in.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("ReviewService.Create: %w", err)
	}
	if err := s.projects.UpdateRatingStats(ctx, in.ProjectID, stats); err != nil {
		return nil, fmt.Errorf("ReviewService.Create: %w", err)
	}

	if s.scoreEnq != nil {
		_ = s.scoreEnq.EnqueueScoreRecalc(ctx, project.BuilderID)
	}
	return r, nil
}

// ListByProject returns paginated reviews for a project.
func (s *ReviewService) ListByProject(ctx context.Context, projectID string, page, limit int) ([]*domain.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	out, total, err := s.reviews.ListByProject(ctx, projectID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("ReviewService.ListByProject: %w", err)
	}
	return out, total, nil
}
