// Package conv provides safe type conversion helpers.
package conv

import (
	"strconv"
	"strings"
)

// StringToInt converts s to int, returning fallback on parse error.
func StringToInt(s string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return fallback
	}
	return n
}

// StringToFloat64 converts s to float64, returning fallback on parse error.
func StringToFloat64(s string, fallback float64) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return fallback
	}
	return f
}

// StringToBool converts s to bool ("true","1","yes" = true), returning fallback on unknown.
func StringToBool(s string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	}
	return fallback
}

// IntToString converts n to its decimal string representation.
func IntToString(n int) string {
	return strconv.Itoa(n)
}

// Float64ToString converts f to a string with up to 2 decimal places.
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// Clamp returns v clamped to [min, max].
func Clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ClampFloat64 returns v clamped to [min, max].
func ClampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// BoolToInt returns 1 if b is true, 0 otherwise.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
