package tags_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/tags"
)

func TestNormalize_DeduplicatesAndLowercases(t *testing.T) {
	input := []string{"Go", "go", " Python ", "rust", "RUST"}
	got := tags.Normalize(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 unique tags, got %d: %v", len(got), got)
	}
	for _, tag := range got {
		if tag != "go" && tag != "python" && tag != "rust" {
			t.Errorf("unexpected tag %q", tag)
		}
	}
}

func TestNormalize_RemovesEmpty(t *testing.T) {
	input := []string{"", "  ", "go", ""}
	got := tags.Normalize(input)
	if len(got) != 1 || got[0] != "go" {
		t.Errorf("expected [go], got %v", got)
	}
}

func TestValidate_TooManyTags(t *testing.T) {
	input := make([]string, 11)
	for i := range input {
		input[i] = "tag"
	}
	if msg := tags.Validate(input); msg == "" {
		t.Error("expected validation error for >10 tags")
	}
}

func TestValidate_TagTooLong(t *testing.T) {
	long := make([]byte, 31)
	for i := range long {
		long[i] = 'a'
	}
	if msg := tags.Validate([]string{string(long)}); msg == "" {
		t.Error("expected validation error for tag >30 chars")
	}
}

func TestValidate_ValidTags(t *testing.T) {
	input := []string{"go", "rust", "typescript", "saas"}
	if msg := tags.Validate(input); msg != "" {
		t.Errorf("expected valid tags to pass, got: %s", msg)
	}
}

func TestContains_CaseInsensitive(t *testing.T) {
	list := []string{"Go", "Rust", "TypeScript"}
	if !tags.Contains(list, "go") {
		t.Error("expected Contains to find 'go' case-insensitively")
	}
	if !tags.Contains(list, "TYPESCRIPT") {
		t.Error("expected Contains to find 'TYPESCRIPT'")
	}
	if tags.Contains(list, "python") {
		t.Error("expected Contains to NOT find 'python'")
	}
}
