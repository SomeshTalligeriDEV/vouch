package timeutil_test

import (
	"testing"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/timeutil"
)

func TestStartOfDay(t *testing.T) {
	ts := time.Date(2024, 6, 15, 14, 30, 45, 999, time.UTC)
	got := timeutil.StartOfDay(ts)
	want := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("StartOfDay: expected %v, got %v", want, got)
	}
}

func TestEndOfDay(t *testing.T) {
	ts := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	got := timeutil.EndOfDay(ts)
	if got.Hour() != 23 || got.Minute() != 59 || got.Second() != 59 {
		t.Errorf("EndOfDay: unexpected time %v", got)
	}
}

func TestDaysBetween_SameDay(t *testing.T) {
	a := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	b := time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC)
	if d := timeutil.DaysBetween(a, b); d != 0 {
		t.Errorf("expected 0 days, got %d", d)
	}
}

func TestDaysBetween_OneDayApart(t *testing.T) {
	a := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	if d := timeutil.DaysBetween(a, b); d != 1 {
		t.Errorf("expected 1 day, got %d", d)
	}
}

func TestDaysBetween_Reversed(t *testing.T) {
	a := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	b := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if d := timeutil.DaysBetween(a, b); d != 4 {
		t.Errorf("expected 4 days (reversed), got %d", d)
	}
}

func TestIsWithinDays(t *testing.T) {
	recent := time.Now().UTC().AddDate(0, 0, -3)
	if !timeutil.IsWithinDays(recent, 7) {
		t.Error("expected 3 days ago to be within 7 days")
	}
	old := time.Now().UTC().AddDate(0, 0, -30)
	if timeutil.IsWithinDays(old, 7) {
		t.Error("expected 30 days ago to NOT be within 7 days")
	}
}
