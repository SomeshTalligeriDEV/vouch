package validate_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validate"
)

func TestUsername(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"alice", true},
		{"Alice_123", true},
		{"a-b", true},
		{"ab", false},        // too short
		{"", false},          // empty
		{"has space", false}, // space not allowed
		{"toolongusernamethatexceedsthirtytwocharacters", false},
	}
	for _, c := range cases {
		got := validate.Username(c.in)
		if got != c.want {
			t.Errorf("Username(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestPassword(t *testing.T) {
	if validate.Password("short") {
		t.Error("expected false for short password")
	}
	if !validate.Password("longenough") {
		t.Error("expected true for 10-char password")
	}
	if !validate.Password("exactly8") {
		t.Error("expected true for 8-char password")
	}
}

func TestURL(t *testing.T) {
	valid := []string{
		"https://example.com",
		"http://localhost:3000",
		"https://sub.domain.io/path?q=1",
	}
	for _, u := range valid {
		if !validate.URL(u) {
			t.Errorf("expected URL(%q) = true", u)
		}
	}

	invalid := []string{
		"",
		"javascript:alert(1)",
		"ftp://files.example.com",
		"not-a-url",
		"//example.com",
	}
	for _, u := range invalid {
		if validate.URL(u) {
			t.Errorf("expected URL(%q) = false", u)
		}
	}
}

func TestNonEmpty(t *testing.T) {
	if validate.NonEmpty("") {
		t.Error("expected false for empty")
	}
	if validate.NonEmpty("   ") {
		t.Error("expected false for whitespace-only")
	}
	if !validate.NonEmpty("hi") {
		t.Error("expected true for non-empty")
	}
}

func TestMaxLen(t *testing.T) {
	if !validate.MaxLen("hello", 5) {
		t.Error("expected true for len==max")
	}
	if validate.MaxLen("hello", 4) {
		t.Error("expected false when over max")
	}
}

func TestMinLen(t *testing.T) {
	if !validate.MinLen("hello", 5) {
		t.Error("expected true for len==min")
	}
	if validate.MinLen("hi", 5) {
		t.Error("expected false when under min")
	}
}

func TestInRange(t *testing.T) {
	if !validate.InRange(5, 1, 10) {
		t.Error("expected true for value in range")
	}
	if validate.InRange(0, 1, 10) {
		t.Error("expected false for value below range")
	}
	if validate.InRange(11, 1, 10) {
		t.Error("expected false for value above range")
	}
}

func TestOneOf(t *testing.T) {
	if !validate.OneOf("gold", "bronze", "silver", "gold") {
		t.Error("expected true for allowed value")
	}
	if validate.OneOf("platinum", "bronze", "silver", "gold") {
		t.Error("expected false for disallowed value")
	}
}
