package strutil_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/strutil"
)

func TestTruncate_Short(t *testing.T) {
	got := strutil.Truncate("hello", "...", 10)
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestTruncate_Long(t *testing.T) {
	got := strutil.Truncate("hello world", "...", 5)
	if got != "hello..." {
		t.Errorf("expected 'hello...', got %q", got)
	}
}

func TestToTitleCase(t *testing.T) {
	got := strutil.ToTitleCase("hello world")
	if got != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", got)
	}
}

func TestCountWords(t *testing.T) {
	if strutil.CountWords("one two three") != 3 {
		t.Error("expected 3 words")
	}
	if strutil.CountWords("  ") != 0 {
		t.Error("expected 0 for whitespace")
	}
}

func TestIsBlank(t *testing.T) {
	if !strutil.IsBlank("") {
		t.Error("expected true for empty")
	}
	if !strutil.IsBlank("   ") {
		t.Error("expected true for whitespace")
	}
	if strutil.IsBlank("a") {
		t.Error("expected false for 'a'")
	}
}

func TestCamelToSnake(t *testing.T) {
	cases := map[string]string{
		"camelCase":      "camel_case",
		"BuilderScore":   "builder_score",
		"totalScore":     "total_score",
		"lowercase":      "lowercase",
	}
	for in, want := range cases {
		got := strutil.CamelToSnake(in)
		if got != want {
			t.Errorf("CamelToSnake(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRepeat(t *testing.T) {
	got := strutil.Repeat("ab", "-", 3)
	if got != "ab-ab-ab" {
		t.Errorf("expected 'ab-ab-ab', got %q", got)
	}
	if strutil.Repeat("x", ",", 0) != "" {
		t.Error("expected empty for n=0")
	}
}

func TestContains(t *testing.T) {
	if !strutil.Contains([]string{"a", "b", "c"}, "b") {
		t.Error("expected true")
	}
	if strutil.Contains([]string{"a", "b"}, "B") {
		t.Error("expected false (case-sensitive)")
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	if !strutil.ContainsIgnoreCase([]string{"Gold", "Silver"}, "gold") {
		t.Error("expected true")
	}
	if strutil.ContainsIgnoreCase([]string{"Gold"}, "platinum") {
		t.Error("expected false")
	}
}
