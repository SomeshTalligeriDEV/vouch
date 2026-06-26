// Package retry provides simple retry logic with exponential backoff.
package retry

import (
	"context"
	"errors"
	"time"
)

// Config holds retry parameters.
type Config struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
	Multiplier  float64
}

// Default returns a sensible default Config (3 attempts, 100ms initial wait).
func Default() Config {
	return Config{
		MaxAttempts: 3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     5 * time.Second,
		Multiplier:  2.0,
	}
}

// ErrMaxAttemptsReached is returned when all attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("max retry attempts reached")

// Do runs fn up to cfg.MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns nil.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = 2.0
	}

	wait := cfg.InitialWait
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}

		wait = time.Duration(float64(wait) * cfg.Multiplier)
		if cfg.MaxWait > 0 && wait > cfg.MaxWait {
			wait = cfg.MaxWait
		}
	}

	return lastErr
}

// DoSimple runs fn up to maxAttempts times with no backoff.
func DoSimple(ctx context.Context, maxAttempts int, fn func() error) error {
	return Do(ctx, Config{MaxAttempts: maxAttempts, InitialWait: 0}, fn)
}
