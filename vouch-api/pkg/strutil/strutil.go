// Package strutil provides common string utility functions.
package strutil

import (
	"strings"
	"unicode"
)

// Truncate shortens s to at most n runes, appending suffix if truncated.
func Truncate(s, suffix string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + suffix
}

// ToTitleCase converts each word's first letter to uppercase.
func ToTitleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(string([]rune(w)[:1])) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// CountWords returns the number of whitespace-separated words in s.
func CountWords(s string) int {
	return len(strings.Fields(s))
}

// IsBlank returns true if s contains only whitespace or is empty.
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// CamelToSnake converts camelCase to snake_case.
func CamelToSnake(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			sb.WriteByte('_')
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}

// Repeat returns s repeated n times with sep between repetitions.
func Repeat(s, sep string, n int) string {
	if n <= 0 {
		return ""
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = s
	}
	return strings.Join(parts, sep)
}

// Contains reports whether any element of slice equals target (case-sensitive).
func Contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// ContainsIgnoreCase reports whether any element of slice equals target (case-insensitive).
func ContainsIgnoreCase(slice []string, target string) bool {
	lower := strings.ToLower(target)
	for _, s := range slice {
		if strings.ToLower(s) == lower {
			return true
		}
	}
	return false
}
