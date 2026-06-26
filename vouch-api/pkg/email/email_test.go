package email_test

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/email"
)

func TestValidate_ValidAddresses(t *testing.T) {
	valid := []string{
		"user@example.com",
		"user+tag@sub.domain.org",
		"first.last@company.io",
		"alice123@test.co.uk",
	}
	for _, addr := range valid {
		if err := email.Validate(addr); err != nil {
			t.Errorf("Validate(%q): expected valid, got error: %v", addr, err)
		}
	}
}

func TestValidate_InvalidAddresses(t *testing.T) {
	invalid := []string{
		"",
		"notanemail",
		"@nodomain.com",
		"user@",
		"user@domain",
		"spaces in@example.com",
		"user@@double.com",
	}
	for _, addr := range invalid {
		if err := email.Validate(addr); !errors.Is(err, email.ErrInvalidEmail) {
			t.Errorf("Validate(%q): expected ErrInvalidEmail, got %v", addr, err)
		}
	}
}

func TestNormalize_Lowercases(t *testing.T) {
	got := email.Normalize("USER@EXAMPLE.COM")
	if got != "user@example.com" {
		t.Errorf("expected 'user@example.com', got %q", got)
	}
}

func TestNormalize_TrimsWhitespace(t *testing.T) {
	got := email.Normalize("  user@example.com  ")
	if got != "user@example.com" {
		t.Errorf("expected 'user@example.com', got %q", got)
	}
}

func TestNormalizeAndValidate_Valid(t *testing.T) {
	addr, err := email.NormalizeAndValidate("  USER@EXAMPLE.COM  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "user@example.com" {
		t.Errorf("expected 'user@example.com', got %q", addr)
	}
}

func TestNormalizeAndValidate_Invalid(t *testing.T) {
	_, err := email.NormalizeAndValidate("notvalid")
	if !errors.Is(err, email.ErrInvalidEmail) {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}
