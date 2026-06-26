package config_test

import (
	"os"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/config"
)

func setEnv(t *testing.T, key, val string) {
	t.Helper()
	t.Setenv(key, val)
}

func baseEnv(t *testing.T) {
	t.Helper()
	setEnv(t, "MONGO_URI", "mongodb://localhost:27017")
	setEnv(t, "MONGO_DB", "vouch_test")
	setEnv(t, "REDIS_URL", "redis://localhost:6379")
	setEnv(t, "JWT_SECRET", "a-secure-32-character-dev-secret!")
	setEnv(t, "JWT_REFRESH_SECRET", "another-32-char-dev-secret-here!")
	setEnv(t, "GITHUB_CLIENT_ID", "gh_id")
	setEnv(t, "GITHUB_CLIENT_SECRET", "gh_secret")
}

func TestConfig_ValidDev(t *testing.T) {
	baseEnv(t)
	setEnv(t, "ENV", "development")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Env != "development" {
		t.Fatalf("expected env 'development', got %q", cfg.Env)
	}
}

func TestConfig_MissingRequiredVar(t *testing.T) {
	for _, key := range []string{"MONGO_URI", "REDIS_URL", "JWT_SECRET", "JWT_REFRESH_SECRET"} {
		t.Run("missing_"+key, func(t *testing.T) {
			baseEnv(t)
			os.Unsetenv(key)
			_, err := config.Load()
			if err == nil {
				t.Fatalf("expected error when %s is missing", key)
			}
		})
	}
}

func TestConfig_ProductionRejectsWeakSecret(t *testing.T) {
	baseEnv(t)
	setEnv(t, "ENV", "production")
	setEnv(t, "JWT_SECRET", "change-me-in-production")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for weak JWT_SECRET in production")
	}
}

func TestConfig_ProductionRejectsShortSecret(t *testing.T) {
	baseEnv(t)
	setEnv(t, "ENV", "production")
	setEnv(t, "JWT_SECRET", "too-short")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for short JWT_SECRET in production")
	}
}

func TestConfig_ProductionAcceptsStrongSecret(t *testing.T) {
	baseEnv(t)
	setEnv(t, "ENV", "production")
	setEnv(t, "JWT_SECRET", "a-very-long-and-secure-production-jwt-secret-value!")
	setEnv(t, "JWT_REFRESH_SECRET", "another-very-long-and-secure-production-refresh-secret!")
	_, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error with strong secrets: %v", err)
	}
}
