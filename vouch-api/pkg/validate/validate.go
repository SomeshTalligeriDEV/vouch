// Package validate provides lightweight input validation helpers.
package validate

import (
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

// Username returns true if s is a valid username (3-32 chars, alphanumeric/_/-).
func Username(s string) bool {
	return usernameRe.MatchString(s)
}

// Password returns true if s meets minimum password requirements (8+ chars).
func Password(s string) bool {
	return utf8.RuneCountInString(s) >= 8
}

// URL returns true if s is a valid absolute HTTP or HTTPS URL.
func URL(s string) bool {
	if s == "" {
		return false
	}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// NonEmpty returns true if s is not empty after trimming whitespace.
func NonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// MaxLen returns true if s has at most n UTF-8 characters.
func MaxLen(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

// MinLen returns true if s has at least n UTF-8 characters.
func MinLen(s string, n int) bool {
	return utf8.RuneCountInString(s) >= n
}

// InRange returns true if n is between min and max inclusive.
func InRange(n, min, max int) bool {
	return n >= min && n <= max
}

// OneOf returns true if s matches one of the allowed values.
func OneOf(s string, allowed ...string) bool {
	for _, a := range allowed {
		if s == a {
			return true
		}
	}
	return false
}
