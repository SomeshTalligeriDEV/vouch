package sanitize_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/sanitize"
)

func TestText_TrimsAndLimits(t *testing.T) {
	input := "  hello world  "
	got := sanitize.Text(input, 5)
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestText_RemovesControlChars(t *testing.T) {
	input := "hello\x00world\x07"
	got := sanitize.Text(input, 0)
	if got != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", got)
	}
}

func TestText_ZeroMaxLen_NoTruncation(t *testing.T) {
	input := "a very long string that should not be truncated"
	got := sanitize.Text(input, 0)
	if got != input {
		t.Errorf("expected unchanged string, got %q", got)
	}
}

func TestText_PreservesNewlines(t *testing.T) {
	input := "line1\nline2\r\nline3"
	got := sanitize.Text(input, 0)
	if got != input {
		t.Errorf("expected newlines preserved, got %q", got)
	}
}

func TestHTML_EscapesAngleBrackets(t *testing.T) {
	got := sanitize.HTML("<script>alert('xss')</script>")
	if got == "<script>alert('xss')</script>" {
		t.Error("expected HTML to be escaped")
	}
	if got != "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;" {
		t.Errorf("unexpected escape result: %q", got)
	}
}

func TestSlug_BasicConversion(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Hello World!", "hello-world"},
		{"  foo   bar  ", "foo-bar"},
		{"already-ok", "already-ok"},
		{"UPPERCASE", "uppercase"},
		{"", ""},
	}
	for _, tc := range cases {
		got := sanitize.Slug(tc.in)
		if got != tc.want {
			t.Errorf("Slug(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestUsername_FiltersInvalidChars(t *testing.T) {
	got := sanitize.Username("Alice_Smith-99@example.com")
	if got != "alice_smith-99examplecom" {
		t.Errorf("Username: unexpected result %q", got)
	}
}

func TestUsername_EmptyInput(t *testing.T) {
	got := sanitize.Username("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
