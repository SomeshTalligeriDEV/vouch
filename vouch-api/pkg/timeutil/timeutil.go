// Package timeutil provides time-related helpers.
package timeutil

import (
	"time"
)

// StartOfDay returns the start of the UTC day for t.
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// EndOfDay returns the end of the UTC day for t (23:59:59.999999999).
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).Add(24*time.Hour - 1)
}

// DaysBetween returns the number of whole days between a and b (ignoring time of day).
func DaysBetween(a, b time.Time) int {
	aDay := StartOfDay(a)
	bDay := StartOfDay(b)
	diff := bDay.Sub(aDay)
	if diff < 0 {
		diff = -diff
	}
	return int(diff.Hours() / 24)
}

// IsToday returns true if t is today in UTC.
func IsToday(t time.Time) bool {
	return DaysBetween(t, time.Now().UTC()) == 0
}

// IsWithinDays returns true if t is within n days of now.
func IsWithinDays(t time.Time, n int) bool {
	return DaysBetween(t, time.Now().UTC()) <= n
}
