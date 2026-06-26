package search_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/search"
)

func TestNormalize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Hello World!", "hello world"},
		{"TypeScript 5.x", "typescript 5x"},
		{"  spaces  ", "  spaces  "},
		{"", ""},
	}
	for _, tc := range cases {
		got := search.Normalize(tc.in)
		if got != tc.want {
			t.Errorf("Normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestContains_Match(t *testing.T) {
	if !search.Contains("Hello World", "world") {
		t.Error("expected case-insensitive match")
	}
	if !search.Contains("TypeScript project", "type") {
		t.Error("expected prefix match")
	}
}

func TestContains_NoMatch(t *testing.T) {
	if search.Contains("Hello", "xyz") {
		t.Error("expected no match")
	}
}

func TestContains_EmptyQuery(t *testing.T) {
	if !search.Contains("anything", "") {
		t.Error("empty query should always match")
	}
}

func TestScore_ExactMatch(t *testing.T) {
	s := search.Score("vouch", "vouch")
	if s != 1.0 {
		t.Errorf("exact match score should be 1.0, got %f", s)
	}
}

func TestScore_NoMatch(t *testing.T) {
	s := search.Score("vouch", "xyz")
	if s != 0 {
		t.Errorf("no match score should be 0, got %f", s)
	}
}

func TestScore_PartialMatch(t *testing.T) {
	s := search.Score("vouch platform", "vouch")
	if s <= 0.5 {
		t.Errorf("partial prefix match should score > 0.5, got %f", s)
	}
}

func TestFilterSlice(t *testing.T) {
	items := []string{"Go project", "Rust project", "TypeScript app", "Go server"}
	got := search.FilterSlice(items, "go")
	if len(got) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(got), got)
	}
}

func TestFilterSlice_EmptyQuery(t *testing.T) {
	items := []string{"a", "b", "c"}
	got := search.FilterSlice(items, "")
	if len(got) != 3 {
		t.Errorf("empty query should return all items, got %d", len(got))
	}
}
