// Package sanitize provides input sanitization helpers.
package sanitize

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

var controlChars = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

// Text removes control characters (excluding \t, \n, \r), trims whitespace,
// and limits the result to maxLen runes.
func Text(s string, maxLen int) string {
	s = controlChars.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	if maxLen > 0 && utf8.RuneCountInString(s) > maxLen {
		runes := []rune(s)
		s = string(runes[:maxLen])
	}
	return s
}

// HTML escapes HTML special characters to prevent XSS when embedding in HTML.
func HTML(s string) string {
	return html.EscapeString(s)
}

// Slug normalizes a string to be safe for use as a URL slug.
// Non-alphanumeric characters become hyphens; consecutive hyphens collapse.
var slugNonWord = regexp.MustCompile(`[^a-z0-9]+`)

func Slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugNonWord.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// Username removes characters that are not valid in a username.
// Only a-z, 0-9, underscore and hyphen are allowed.
var nonUsername = regexp.MustCompile(`[^a-z0-9_-]`)

func Username(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return nonUsername.ReplaceAllString(s, "")
}
