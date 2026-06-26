package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

// GitHubGateway abstracts the GitHub OAuth + user API. It exchanges an OAuth
// code for an access token and reads the authenticated user's profile.
type GitHubGateway interface {
	ExchangeCode(ctx context.Context, code string) (GitHubProfile, error)
}

// UserService holds business logic for users and authentication.
type UserService struct {
	users  domain.UserRepository
	jwt    *jwt.Manager
	github GitHubGateway
}

// NewUserService constructs a UserService.
func NewUserService(users domain.UserRepository, jwtMgr *jwt.Manager, github GitHubGateway) *UserService {
	return &UserService{users: users, jwt: jwtMgr, github: github}
}

// GitHubProfile is the subset of GitHub user data Vouch consumes.
type GitHubProfile struct {
	ID        int64
	Login     string
	Name      string
	Email     string
	AvatarURL string
	Bio       string
}

// LoginWithGitHub exchanges an OAuth code for a profile, then finds-or-creates
// the corresponding user and issues tokens.
func (s *UserService) LoginWithGitHub(ctx context.Context, code string) (*domain.User, *jwt.TokenPair, error) {
	profile, err := s.github.ExchangeCode(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("UserService.LoginWithGitHub: %w", err)
	}
	return s.UpsertFromGitHub(ctx, profile)
}

// UpsertFromGitHub finds-or-creates a user from a GitHub profile and returns a
// fresh token pair. This is the heart of the GitHub OAuth login flow.
func (s *UserService) UpsertFromGitHub(ctx context.Context, p GitHubProfile) (*domain.User, *jwt.TokenPair, error) {
	user, err := s.users.GetByGitHubID(ctx, p.ID)
	switch {
	case err == nil:
		// Existing user — nothing to create.
	case isNotFound(err):
		user = &domain.User{
			Email:       p.Email,
			Username:    s.uniqueUsername(ctx, p.Login),
			Name:        p.Name,
			Bio:         p.Bio,
			AvatarURL:   p.AvatarURL,
			GitHubID:    p.ID,
			GitHubLogin: p.Login,
			Role:        domain.RoleBuilder,
			IsVerified:  true,
		}
		if err := s.users.Create(ctx, user); err != nil {
			return nil, nil, fmt.Errorf("UserService.UpsertFromGitHub: %w", err)
		}
	default:
		return nil, nil, fmt.Errorf("UserService.UpsertFromGitHub: %w", err)
	}

	tokens, err := s.jwt.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, nil, fmt.Errorf("UserService.UpsertFromGitHub: %w", err)
	}
	return user, tokens, nil
}

// uniqueUsername returns login if free, otherwise appends a short suffix.
func (s *UserService) uniqueUsername(ctx context.Context, login string) string {
	login = strings.ToLower(strings.TrimSpace(login))
	if login == "" {
		login = "builder"
	}
	if _, err := s.users.GetByUsername(ctx, login); isNotFound(err) {
		return login
	}
	return login + "-" + randHex(2)
}

// Refresh validates a refresh token and issues a new token pair.
func (s *UserService) Refresh(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	claims, err := s.jwt.VerifyRefresh(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("UserService.Refresh: %w", domain.ErrUnauthorized)
	}
	user, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("UserService.Refresh: %w", err)
	}
	tokens, err := s.jwt.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("UserService.Refresh: %w", err)
	}
	return tokens, nil
}

// GetByUsername returns a public profile by username.
func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	u, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetByUsername: %w", err)
	}
	return u, nil
}

// GetByID returns a user by id.
func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetByID: %w", err)
	}
	return u, nil
}

// UpdateInput carries the editable fields of a profile.
type UpdateInput struct {
	Name          string
	Bio           string
	AvatarURL     string
	WebsiteURL    string
	TwitterHandle string
}

// UpdateProfile applies editable fields to the authenticated user's profile.
func (s *UserService) UpdateProfile(ctx context.Context, userID string, in UpdateInput) (*domain.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UserService.UpdateProfile: %w", err)
	}
	u.Name = in.Name
	u.Bio = in.Bio
	if in.AvatarURL != "" {
		u.AvatarURL = in.AvatarURL
	}
	u.WebsiteURL = in.WebsiteURL
	u.TwitterHandle = in.TwitterHandle
	if err := s.users.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("UserService.UpdateProfile: %w", err)
	}
	return u, nil
}
