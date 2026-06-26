package domain_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestUser_DisplayName_UsesNameWhenSet(t *testing.T) {
	u := domain.User{Username: "alice", Name: "Alice Smith"}
	if u.DisplayName() != "Alice Smith" {
		t.Fatalf("expected 'Alice Smith', got %q", u.DisplayName())
	}
}

func TestUser_DisplayName_FallsBackToUsername(t *testing.T) {
	u := domain.User{Username: "alice"}
	if u.DisplayName() != "alice" {
		t.Fatalf("expected 'alice', got %q", u.DisplayName())
	}
}

func TestUser_IsAdmin_TrueForAdminRole(t *testing.T) {
	u := domain.User{Role: "admin"}
	if !u.IsAdmin() {
		t.Fatal("expected IsAdmin() true for role 'admin'")
	}
}

func TestUser_IsAdmin_FalseForUserRole(t *testing.T) {
	u := domain.User{Role: "user"}
	if u.IsAdmin() {
		t.Fatal("expected IsAdmin() false for role 'user'")
	}
}
