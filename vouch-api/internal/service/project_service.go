package service

import (
	"context"
	"fmt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// ProjectService holds business logic for projects.
type ProjectService struct {
	projects domain.ProjectRepository
	scoreEnq ScoreEnqueuer
}

// NewProjectService constructs a ProjectService.
func NewProjectService(projects domain.ProjectRepository, scoreEnq ScoreEnqueuer) *ProjectService {
	return &ProjectService{projects: projects, scoreEnq: scoreEnq}
}

// CreateInput carries the fields needed to create a project.
type CreateInput struct {
	Title       string
	Tagline     string
	Description string
	LogoURL     string
	LiveURL     string
	RepoURL     string
	PaymentLink string
	Tags        []string
	ForSale     bool
	AskPrice    float64
}

// Create creates a new project owned by the given builder.
func (s *ProjectService) Create(ctx context.Context, builderID string, in CreateInput) (*domain.Project, error) {
	p := &domain.Project{
		BuilderID:   builderID,
		Title:       in.Title,
		Slug:        slugify(in.Title),
		Tagline:     in.Tagline,
		Description: in.Description,
		LogoURL:     in.LogoURL,
		LiveURL:     in.LiveURL,
		RepoURL:     in.RepoURL,
		PaymentLink: in.PaymentLink,
		Tags:        in.Tags,
		Status:      domain.ProjectStatusLive,
		ForSale:     in.ForSale,
		AskPrice:    in.AskPrice,
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}
	if err := s.projects.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("ProjectService.Create: %w", err)
	}
	s.scheduleRecalc(ctx, builderID)
	return p, nil
}

// Get returns a single project by id.
func (s *ProjectService) Get(ctx context.Context, id string) (*domain.Project, error) {
	p, err := s.projects.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.Get: %w", err)
	}
	return p, nil
}

// List returns a filtered, paginated set of projects.
func (s *ProjectService) List(ctx context.Context, f domain.ProjectFilter) ([]*domain.Project, int64, error) {
	normalizePage(&f.Page, &f.Limit)
	out, total, err := s.projects.List(ctx, f)
	if err != nil {
		return nil, 0, fmt.Errorf("ProjectService.List: %w", err)
	}
	return out, total, nil
}

// ListByBuilder returns all of a builder's projects.
func (s *ProjectService) ListByBuilder(ctx context.Context, builderID string) ([]*domain.Project, error) {
	out, err := s.projects.ListByBuilder(ctx, builderID)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.ListByBuilder: %w", err)
	}
	return out, nil
}

// UpdateInput carries the editable fields of a project.
type ProjectUpdateInput struct {
	Title       string
	Tagline     string
	Description string
	LogoURL     string
	LiveURL     string
	RepoURL     string
	PaymentLink string
	Tags        []string
	ForSale     bool
	AskPrice    float64
	Status      domain.ProjectStatus
}

// Update mutates a project after verifying the caller owns it.
func (s *ProjectService) Update(ctx context.Context, id, builderID string, in ProjectUpdateInput) (*domain.Project, error) {
	p, err := s.ownedProject(ctx, id, builderID)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.Update: %w", err)
	}
	if in.Title != "" {
		p.Title = in.Title
	}
	p.Tagline = in.Tagline
	p.Description = in.Description
	p.LogoURL = in.LogoURL
	p.LiveURL = in.LiveURL
	p.RepoURL = in.RepoURL
	p.PaymentLink = in.PaymentLink
	if in.Tags != nil {
		p.Tags = in.Tags
	}
	p.ForSale = in.ForSale
	p.AskPrice = in.AskPrice
	if in.Status != "" {
		p.Status = in.Status
	}
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("ProjectService.Update: %w", err)
	}
	return p, nil
}

// Archive soft-deletes a project (sets status to archived).
func (s *ProjectService) Archive(ctx context.Context, id, builderID string) error {
	p, err := s.ownedProject(ctx, id, builderID)
	if err != nil {
		return fmt.Errorf("ProjectService.Archive: %w", err)
	}
	p.Status = domain.ProjectStatusArchived
	if err := s.projects.Update(ctx, p); err != nil {
		return fmt.Errorf("ProjectService.Archive: %w", err)
	}
	s.scheduleRecalc(ctx, builderID)
	return nil
}

func (s *ProjectService) ownedProject(ctx context.Context, id, builderID string) (*domain.Project, error) {
	p, err := s.projects.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.BuilderID != builderID {
		return nil, domain.ErrForbidden
	}
	return p, nil
}

func (s *ProjectService) scheduleRecalc(ctx context.Context, builderID string) {
	if s.scoreEnq == nil {
		return
	}
	// Best-effort: a failed enqueue should not fail the user's request.
	_ = s.scoreEnq.EnqueueScoreRecalc(ctx, builderID)
}

// normalizePage clamps pagination inputs to sane defaults.
func normalizePage(page, limit *int) {
	if *page < 1 {
		*page = 1
	}
	if *limit <= 0 || *limit > 100 {
		*limit = 20
	}
}
