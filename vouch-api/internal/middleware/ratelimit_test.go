package middleware_test

import (
	"testing"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
)

func TestNewRateLimiter_StoresConfig(t *testing.T) {
	rl := middleware.NewRateLimiter(nil, 100, time.Minute)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
}

func TestNewRateLimiter_LimitIsPositive(t *testing.T) {
	rl := middleware.NewRateLimiter(nil, 50, 30*time.Second)
	if rl == nil {
		t.Fatal("expected non-nil limiter")
	}
	h := rl.Limit()
	if h == nil {
		t.Fatal("Limit() returned nil handler")
	}
}
