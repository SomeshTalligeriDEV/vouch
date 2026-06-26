package retry_test

import (
	"context"
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/retry"
)

var errFake = errors.New("fake error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 3}, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 3, InitialWait: 0}, func() error {
		calls++
		if calls < 3 {
			return errFake
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 3, InitialWait: 0}, func() error {
		calls++
		return errFake
	})
	if err == nil {
		t.Error("expected error after exhausting attempts")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, retry.Config{MaxAttempts: 3, InitialWait: 0}, func() error {
		calls++
		return errFake
	})
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestDoSimple_SucceedsOnSecondAttempt(t *testing.T) {
	calls := 0
	err := retry.DoSimple(context.Background(), 5, func() error {
		calls++
		if calls < 2 {
			return errFake
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestDefault_ReasonableValues(t *testing.T) {
	cfg := retry.Default()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected 3 attempts, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialWait <= 0 {
		t.Error("expected positive initial wait")
	}
}
