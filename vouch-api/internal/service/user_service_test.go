package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

type fakeGitHub struct{}

func (f *fakeGitHub) ExchangeCode(_ context.Context, _ string) (GitHubProfile, error) {
	return GitHubProfile{}, errors.New("not implemented in tests")
}

func newTestJWTManager() *jwt.Manager {
	return jwt.NewManager(
		"test-secret-32-chars-long-xxxxxxxx",
		"refresh-secret-32-chars-long-xxxx",
	)
}

func TestGetByUsername_NotFound(t *testing.T) {
	svc := NewUserService(newFakeUserRepo(), newTestJWTManager(), &fakeGitHub{})
	_, err := svc.GetByUsername(ctx(), "nobody")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateProfile_SetsFields(t *testing.T) {
	svc := NewUserService(newFakeUserRepo(), newTestJWTManager(), &fakeGitHub{})

	user, _, err := svc.UpsertFromGitHub(ctx(), GitHubProfile{
		ID: 42, Login: "alice", Name: "Alice", Email: "alice@example.com",
	})
	if err != nil {
		t.Fatalf("UpsertFromGitHub: %v", err)
	}

	updated, err := svc.UpdateProfile(ctx(), user.ID, UpdateInput{
		Name: "Alice Updated",
		Bio:  "Building cool stuff",
	})
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Fatalf("expected 'Alice Updated', got %q", updated.Name)
	}
	if updated.Bio != "Building cool stuff" {
		t.Fatalf("expected bio set, got %q", updated.Bio)
	}
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	svc := NewUserService(newFakeUserRepo(), newTestJWTManager(), &fakeGitHub{})
	_, err := svc.UpdateProfile(ctx(), "nonexistent-id", UpdateInput{Name: "X"})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestUpsertFromGitHub_CreatesUser(t *testing.T) {
	svc := NewUserService(newFakeUserRepo(), newTestJWTManager(), &fakeGitHub{})
	user, tokens, err := svc.UpsertFromGitHub(ctx(), GitHubProfile{ID: 99, Login: "bob", Name: "Bob"})
	if err != nil {
		t.Fatalf("UpsertFromGitHub: %v", err)
	}
	if user.Username != "bob" {
		t.Fatalf("expected username 'bob', got %q", user.Username)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("expected non-empty token pair")
	}
}
