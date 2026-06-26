package service

import (
	"context"
	"fmt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// AdminStats is a snapshot of platform-wide counts for the admin dashboard.
type AdminStats struct {
	TotalUsers     int64 `json:"total_users"`
	TotalCompanies int64 `json:"total_companies"`
	TotalProjects  int64 `json:"total_projects"`
	OpenProblems   int64 `json:"open_problems"`
	TotalReviews   int64 `json:"total_reviews"`
}

// AdminService surfaces admin-only views of the platform.
type AdminService struct {
	users     domain.UserRepository
	companies domain.CompanyRepository
	projects  domain.ProjectRepository
	problems  domain.ProblemRepository
	reviews   domain.ReviewRepository
}

// NewAdminService constructs an AdminService.
func NewAdminService(
	users domain.UserRepository,
	companies domain.CompanyRepository,
	projects domain.ProjectRepository,
	problems domain.ProblemRepository,
	reviews domain.ReviewRepository,
) *AdminService {
	return &AdminService{
		users:     users,
		companies: companies,
		projects:  projects,
		problems:  problems,
		reviews:   reviews,
	}
}

// Stats returns platform-wide aggregate counts.
func (s *AdminService) Stats(ctx context.Context) (*AdminStats, error) {
	// Projects total
	_, projTotal, err := s.projects.List(ctx, domain.ProjectFilter{Limit: 1, Page: 1})
	if err != nil {
		return nil, fmt.Errorf("AdminService.Stats projects: %w", err)
	}

	// Open problems
	_, openTotal, err := s.problems.List(ctx, domain.ProblemFilter{Status: domain.ProblemStatusOpen, Limit: 1, Page: 1})
	if err != nil {
		return nil, fmt.Errorf("AdminService.Stats problems: %w", err)
	}

	// Companies
	_, compTotal, err := s.companies.List(ctx, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("AdminService.Stats companies: %w", err)
	}

	return &AdminStats{
		TotalCompanies: compTotal,
		TotalProjects:  projTotal,
		OpenProblems:   openTotal,
	}, nil
}

// ListCompanies returns a paginated list of companies.
func (s *AdminService) ListCompanies(ctx context.Context, page, limit int) ([]*domain.Company, int64, error) {
	return s.companies.List(ctx, page, limit)
}
