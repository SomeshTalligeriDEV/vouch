// Package env provides helpers for reading environment variables with defaults.
package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// String reads an env var, returning fallback if it is empty.
func String(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// MustString reads an env var and panics if it is empty.
func MustString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required environment variable not set: " + key)
	}
	return v
}

// Int reads an env var as an integer, returning fallback on missing/parse error.
func Int(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

// Bool reads an env var as a boolean ("true", "1", "yes" = true).
func Bool(key string, fallback bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch v {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return fallback
	}
}

// Duration reads an env var as a time.Duration string (e.g. "5m", "2h").
func Duration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

// IsProduction returns true if the ENV or ENVIRONMENT variable is "production".
func IsProduction() bool {
	env := strings.ToLower(String("ENV", String("ENVIRONMENT", "development")))
	return env == "production" || env == "prod"
}
