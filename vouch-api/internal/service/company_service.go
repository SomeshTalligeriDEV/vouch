package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	neturl "net/url"

	"golang.org/x/crypto/bcrypt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

// CompanyAuthResponse is returned on register and login.
type CompanyAuthResponse struct {
	Company *domain.Company `json:"company"`
	Tokens  *jwt.TokenPair  `json:"tokens"`
}

// CompanyService handles company registration, authentication, and profile updates.
type CompanyService struct {
	companies domain.CompanyRepository
	jwtMgr    *jwt.Manager
}

// NewCompanyService constructs a CompanyService.
func NewCompanyService(companies domain.CompanyRepository, jwtMgr *jwt.Manager) *CompanyService {
	return &CompanyService{companies: companies, jwtMgr: jwtMgr}
}

// Register creates a new company account and returns tokens.
func (s *CompanyService) Register(ctx context.Context, name, email, password, website, size string) (*CompanyAuthResponse, error) {
	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("CompanyService.Register: %w", domain.ErrInvalidInput)
	}
	if len(password) < 8 {
		return nil, fmt.Errorf("CompanyService.Register: password must be at least 8 characters: %w", domain.ErrInvalidInput)
	}
	if website != "" {
		if err := validateHTTPSURL(website); err != nil {
			return nil, fmt.Errorf("CompanyService.Register: %w", domain.ErrInvalidInput)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.Register: %w", domain.ErrInternal)
	}

	slug := slugify(name) // reuse builder slug helper
	// Ensure slug uniqueness by appending timestamp on collision.
	if _, err := s.companies.GetBySlug(ctx, slug); err == nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix()%10000)
	}

	c := &domain.Company{
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: string(hash),
		Name:         name,
		Slug:         slug,
		Website:      website,
		Size:         domain.CompanySize(size),
	}

	if err := s.companies.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("CompanyService.Register: %w", err)
	}

	tokens, err := s.jwtMgr.GenerateTyped(c.ID, c.Slug, "company", jwt.SubjectTypeCompany)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.Register: %w", err)
	}
	return &CompanyAuthResponse{Company: c, Tokens: tokens}, nil
}

// Login authenticates a company and returns tokens.
func (s *CompanyService) Login(ctx context.Context, email, password string) (*CompanyAuthResponse, error) {
	c, err := s.companies.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		// Return generic error to prevent email enumeration.
		return nil, fmt.Errorf("CompanyService.Login: %w", domain.ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("CompanyService.Login: %w", domain.ErrUnauthorized)
	}

	tokens, err := s.jwtMgr.GenerateTyped(c.ID, c.Slug, "company", jwt.SubjectTypeCompany)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.Login: %w", err)
	}
	return &CompanyAuthResponse{Company: c, Tokens: tokens}, nil
}

// Refresh issues new tokens for a company, revoking the old refresh token (rotation).
func (s *CompanyService) Refresh(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	claims, err := s.jwtMgr.VerifyRefresh(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.Refresh: %w", domain.ErrUnauthorized)
	}
	c, err := s.companies.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.Refresh: %w", domain.ErrUnauthorized)
	}
	// Rotate: revoke the old token before issuing a new one.
	_ = s.jwtMgr.RevokeRefresh(ctx, refreshToken)
	return s.jwtMgr.GenerateTyped(c.ID, c.Slug, "company", jwt.SubjectTypeCompany)
}

// Logout revokes the company's refresh token.
func (s *CompanyService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.jwtMgr.RevokeRefresh(ctx, refreshToken); err != nil {
		return fmt.Errorf("CompanyService.Logout: %w", err)
	}
	return nil
}

// GetByID returns a company by its ID.
func (s *CompanyService) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	return s.companies.GetByID(ctx, id)
}

// GetBySlug returns a company by its slug.
func (s *CompanyService) GetBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	return s.companies.GetBySlug(ctx, slug)
}

// UpdateProfile updates mutable company fields.
func (s *CompanyService) UpdateProfile(ctx context.Context, id string, name, website, logoURL, description, size string) (*domain.Company, error) {
	c, err := s.companies.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.UpdateProfile: %w", err)
	}
	if name != "" {
		c.Name = name
	}
	if website != "" {
		if err := validateHTTPSURL(website); err != nil {
			return nil, fmt.Errorf("CompanyService.UpdateProfile: invalid website URL: %w", domain.ErrInvalidInput)
		}
		c.Website = website
	}
	if logoURL != "" {
		c.LogoURL = logoURL
	}
	if description != "" {
		c.Description = description
	}
	if size != "" {
		c.Size = domain.CompanySize(size)
	}
	if err := s.companies.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("CompanyService.UpdateProfile: %w", err)
	}
	return c, nil
}

// validateHTTPSURL rejects any URL whose scheme is not http or https,
// preventing javascript: and data: URIs from being stored and later
// rendered as hrefs.
func validateHTTPSURL(raw string) error {
	u, err := neturl.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		return nil
	default:
		return fmt.Errorf("URL scheme %q not allowed; must be http or https", u.Scheme)
	}
}
