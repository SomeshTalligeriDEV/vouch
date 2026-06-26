// Package tags provides tag normalization and validation.
package tags

import (
	"strings"
	"unicode/utf8"
)

const (
	MaxTags    = 10
	MaxTagLen  = 30
	MinTagLen  = 1
)

// Normalize lowercases and trims whitespace from each tag,
// deduplicates, and removes empty strings.
func Normalize(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		result = append(result, t)
	}
	return result
}

// Validate returns an error string if the tag list is invalid, or "" if ok.
func Validate(tags []string) string {
	if len(tags) > MaxTags {
		return "too many tags (max 10)"
	}
	for _, t := range tags {
		l := utf8.RuneCountInString(t)
		if l < MinTagLen {
			return "tag too short"
		}
		if l > MaxTagLen {
			return "tag too long (max 30 chars)"
		}
	}
	return ""
}

// Contains reports whether the tag list includes a specific tag (case-insensitive).
func Contains(tags []string, tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, t := range tags {
		if strings.ToLower(t) == tag {
			return true
		}
	}
	return false
}
