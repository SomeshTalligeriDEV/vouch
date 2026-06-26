package validator_test

import (
	"strings"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

type testInput struct {
	Name  string `validate:"required,max=50"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=0,max=150"`
}

func TestValidator_ValidStruct(t *testing.T) {
	v := validator.New()
	input := testInput{Name: "Alice", Email: "alice@example.com", Age: 30}
	if err := v.Struct(input); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidator_MissingRequired(t *testing.T) {
	v := validator.New()
	input := testInput{Email: "alice@example.com"}
	err := v.Struct(input)
	if err == nil {
		t.Fatal("expected validation error for missing Name")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Fatalf("expected error to mention 'name', got: %v", err)
	}
}

func TestValidator_InvalidEmail(t *testing.T) {
	v := validator.New()
	input := testInput{Name: "Alice", Email: "not-an-email", Age: 25}
	err := v.Struct(input)
	if err == nil {
		t.Fatal("expected validation error for invalid email")
	}
	if !strings.Contains(err.Error(), "email") {
		t.Fatalf("expected error to mention 'email', got: %v", err)
	}
}

func TestValidator_MaxLengthExceeded(t *testing.T) {
	v := validator.New()
	longName := strings.Repeat("a", 51)
	input := testInput{Name: longName, Email: "alice@example.com", Age: 25}
	err := v.Struct(input)
	if err == nil {
		t.Fatal("expected validation error for name exceeding max length")
	}
}

func TestValidator_NilSafe(t *testing.T) {
	v := validator.New()
	if v == nil {
		t.Fatal("validator.New() returned nil")
	}
}
