package service

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func newTestCompanyService(t *testing.T) *CompanyService {
	t.Helper()
	return NewCompanyService(newFakeCompanyRepo(), newTestJWT())
}

func TestCompanyService_Register_PasswordTooShort(t *testing.T) {
	svc := newTestCompanyService(t)
	_, err := svc.Register(ctx(), "Acme", "acme@example.com", "abc", "https://acme.com", "small")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for short password, got %v", err)
	}
}

func TestCompanyService_Register_EmptyName(t *testing.T) {
	svc := newTestCompanyService(t)
	_, err := svc.Register(ctx(), "", "acme@example.com", "SecurePass123!", "https://acme.com", "small")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for empty name, got %v", err)
	}
}

func TestCompanyService_DuplicateEmailRegistration(t *testing.T) {
	svc := newTestCompanyService(t)
	_, err := svc.Register(ctx(), "Acme", "dup@example.com", "SecurePass123!", "https://acme.com", "small")
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	_, err = svc.Register(ctx(), "Acme2", "dup@example.com", "SecurePass123!", "https://acme2.com", "small")
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists for duplicate email, got %v", err)
	}
}
