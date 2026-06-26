package service

import (
	"strings"
	"testing"
)

func TestSlugify_Basic(t *testing.T) {
	cases := []struct {
		input  string
		prefix string
	}{
		{"Hello World", "hello-world-"},
		{"  leading spaces  ", "leading-spaces-"},
		{"multiple   spaces", "multiple-spaces-"},
		{"Special! Chars@#", "special-chars-"},
		{"already-a-slug", "already-a-slug-"},
		{"CamelCase", "camelcase-"},
		{"numbers123", "numbers123-"},
		{"emoji 🚀 test", "emoji-test-"},
	}

	for _, tc := range cases {
		got := slugify(tc.input)
		if !strings.HasPrefix(got, tc.prefix) {
			t.Errorf("slugify(%q) = %q, want prefix %q", tc.input, got, tc.prefix)
		}
	}
}

func TestSlugify_NeverEmpty(t *testing.T) {
	// A slug from only special characters should not be empty.
	s := slugify("!@#$%^&*()")
	if s == "" {
		t.Error("expected non-empty slug for all-special input")
	}
}

func TestSlugify_HasSuffix(t *testing.T) {
	s := slugify("My Project")
	if !strings.HasPrefix(s, "my-project-") {
		t.Errorf("expected slug to start with 'my-project-', got %q", s)
	}
	// Suffix should be 6 hex chars (3 random bytes = 6 hex chars).
	parts := strings.Split(s, "-")
	suffix := parts[len(parts)-1]
	if len(suffix) != 6 {
		t.Errorf("expected 6-char hex suffix, got %q (len %d)", suffix, len(suffix))
	}
}
