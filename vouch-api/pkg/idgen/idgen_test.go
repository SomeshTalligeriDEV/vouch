package idgen_test

import (
	"strings"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/idgen"
)

func TestHex_Length(t *testing.T) {
	for _, n := range []int{4, 8, 16} {
		got := idgen.Hex(n)
		if len(got) != n*2 {
			t.Errorf("Hex(%d): expected %d chars, got %d", n, n*2, len(got))
		}
	}
}

func TestHex_Unique(t *testing.T) {
	a := idgen.Hex(8)
	b := idgen.Hex(8)
	if a == b {
		t.Error("expected two Hex(8) calls to produce different results")
	}
}

func TestShort_Length(t *testing.T) {
	s := idgen.Short()
	if len(s) != 12 {
		t.Errorf("expected Short() to have 12 chars, got %d", len(s))
	}
}

func TestLong_Length(t *testing.T) {
	l := idgen.Long()
	if len(l) != 32 {
		t.Errorf("expected Long() to have 32 chars, got %d", len(l))
	}
}

func TestPrefixed_Format(t *testing.T) {
	id := idgen.Prefixed("usr")
	if !strings.HasPrefix(id, "usr_") {
		t.Errorf("expected prefix 'usr_', got %s", id)
	}
	if len(id) != 4+16 {
		t.Errorf("expected len=20, got %d: %s", len(id), id)
	}
}

func TestTimeSortable_Format(t *testing.T) {
	id := idgen.TimeSortable()
	// 16 hex chars for timestamp + 16 for random suffix
	if len(id) != 32 {
		t.Errorf("expected TimeSortable() length=32, got %d: %s", len(id), id)
	}
}

func TestIsHex_Valid(t *testing.T) {
	id := idgen.Hex(8)
	if !idgen.IsHex(id, 16) {
		t.Errorf("expected IsHex(%q, 16) = true", id)
	}
}

func TestIsHex_Invalid(t *testing.T) {
	cases := []struct {
		s string
		n int
	}{
		{"abcdefg", 7},  // 'g' not hex
		{"abcd", 6},     // wrong length
		{"ABCDEF12", 8}, // uppercase
	}
	for _, c := range cases {
		if idgen.IsHex(c.s, c.n) {
			t.Errorf("expected IsHex(%q, %d) = false", c.s, c.n)
		}
	}
}
