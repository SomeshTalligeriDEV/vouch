package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, pairs map[string]string) {
	t.Helper()
	for k, v := range pairs {
		t.Setenv(k, v)
	}
}

func baseEnv() map[string]string {
	return map[string]string{
		"MONGO_URI":          "mongodb://localhost:27017",
		"REDIS_URL":          "redis://localhost:6379",
		"JWT_SECRET":         "a-strong-secret-that-is-32-chars!",
		"JWT_REFRESH_SECRET": "another-strong-secret-32-chars!!",
	}
}

func TestLoad_MissingRequiredVars(t *testing.T) {
	// Unset all required vars.
	for _, k := range []string{"MONGO_URI", "REDIS_URL", "JWT_SECRET", "JWT_REFRESH_SECRET"} {
		os.Unsetenv(k)
	}
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when required vars are missing")
	}
}

func TestLoad_Success(t *testing.T) {
	setEnv(t, baseEnv())
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.MongoDB != "vouch" {
		t.Fatalf("expected default mongo db 'vouch', got %q", cfg.MongoDB)
	}
	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
}

func TestLoad_ProductionWeakSecret(t *testing.T) {
	env := baseEnv()
	env["ENV"] = "production"
	env["JWT_SECRET"] = "change-me-in-production"
	setEnv(t, env)
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for placeholder JWT secret in production")
	}
}

func TestLoad_ProductionShortSecret(t *testing.T) {
	env := baseEnv()
	env["ENV"] = "production"
	env["JWT_SECRET"] = "tooshort"
	setEnv(t, env)
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for short JWT secret in production")
	}
}

func TestLoad_ProductionStrongSecret(t *testing.T) {
	env := baseEnv()
	env["ENV"] = "production"
	setEnv(t, env)
	_, err := Load()
	if err != nil {
		t.Fatalf("strong 32-char secret should pass production validation: %v", err)
	}
}
