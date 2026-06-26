package env_test

import (
	"testing"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/env"
)

func TestString_Fallback(t *testing.T) {
	t.Setenv("TEST_STRING_MISSING", "")
	got := env.String("TEST_STRING_MISSING", "default")
	if got != "default" {
		t.Errorf("expected 'default', got %q", got)
	}
}

func TestString_Set(t *testing.T) {
	t.Setenv("TEST_STRING_SET", "hello")
	got := env.String("TEST_STRING_SET", "default")
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestInt_Fallback(t *testing.T) {
	got := env.Int("TEST_INT_MISSING_XYZ123", 42)
	if got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestInt_Set(t *testing.T) {
	t.Setenv("TEST_INT_SET", "99")
	got := env.Int("TEST_INT_SET", 0)
	if got != 99 {
		t.Errorf("expected 99, got %d", got)
	}
}

func TestInt_InvalidValue(t *testing.T) {
	t.Setenv("TEST_INT_INVALID", "notanumber")
	got := env.Int("TEST_INT_INVALID", 5)
	if got != 5 {
		t.Errorf("expected fallback 5, got %d", got)
	}
}

func TestBool_True(t *testing.T) {
	for _, v := range []string{"true", "1", "yes", "TRUE", "YES"} {
		t.Setenv("TEST_BOOL", v)
		if !env.Bool("TEST_BOOL", false) {
			t.Errorf("expected true for value %q", v)
		}
	}
}

func TestBool_False(t *testing.T) {
	for _, v := range []string{"false", "0", "no"} {
		t.Setenv("TEST_BOOL", v)
		if env.Bool("TEST_BOOL", true) {
			t.Errorf("expected false for value %q", v)
		}
	}
}

func TestDuration_Set(t *testing.T) {
	t.Setenv("TEST_DURATION", "5m")
	got := env.Duration("TEST_DURATION", time.Second)
	if got != 5*time.Minute {
		t.Errorf("expected 5m, got %v", got)
	}
}

func TestDuration_Fallback(t *testing.T) {
	got := env.Duration("TEST_DURATION_MISSING_XYZ", 10*time.Second)
	if got != 10*time.Second {
		t.Errorf("expected 10s fallback, got %v", got)
	}
}

func TestIsProduction_Dev(t *testing.T) {
	t.Setenv("ENV", "development")
	if env.IsProduction() {
		t.Error("expected false for development")
	}
}

func TestIsProduction_Prod(t *testing.T) {
	t.Setenv("ENV", "production")
	if !env.IsProduction() {
		t.Error("expected true for production")
	}
}
