package crypto_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/crypto"
)

func TestRandHex_Length(t *testing.T) {
	for _, n := range []int{4, 8, 16, 32} {
		s, err := crypto.RandHex(n)
		if err != nil {
			t.Fatalf("RandHex(%d): %v", n, err)
		}
		if len(s) != n*2 {
			t.Errorf("RandHex(%d): expected len %d, got %d", n, n*2, len(s))
		}
	}
}

func TestRandHex_Uniqueness(t *testing.T) {
	a, _ := crypto.RandHex(16)
	b, _ := crypto.RandHex(16)
	if a == b {
		t.Error("two RandHex(16) calls returned the same value — collision detected")
	}
}

func TestRandBase64URL_NotEmpty(t *testing.T) {
	s, err := crypto.RandBase64URL(32)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) == 0 {
		t.Error("expected non-empty base64url string")
	}
}

func TestConstantTimeEqual_Equal(t *testing.T) {
	if !crypto.ConstantTimeEqual("hello", "hello") {
		t.Error("expected equal strings to return true")
	}
}

func TestConstantTimeEqual_NotEqual(t *testing.T) {
	if crypto.ConstantTimeEqual("hello", "world") {
		t.Error("expected unequal strings to return false")
	}
}

func TestConstantTimeEqual_Empty(t *testing.T) {
	if !crypto.ConstantTimeEqual("", "") {
		t.Error("expected empty strings to be equal")
	}
	if crypto.ConstantTimeEqual("a", "") {
		t.Error("expected non-empty != empty")
	}
}
