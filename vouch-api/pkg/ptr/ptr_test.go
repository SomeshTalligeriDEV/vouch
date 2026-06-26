package ptr_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/ptr"
)

func TestOf_Int(t *testing.T) {
	p := ptr.Of(42)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != 42 {
		t.Errorf("expected 42, got %d", *p)
	}
}

func TestOf_String(t *testing.T) {
	p := ptr.Of("hello")
	if *p != "hello" {
		t.Errorf("expected 'hello', got %q", *p)
	}
}

func TestDeref_NonNil(t *testing.T) {
	n := 99
	if ptr.Deref(&n) != 99 {
		t.Error("expected 99")
	}
}

func TestDeref_Nil(t *testing.T) {
	var p *int
	if ptr.Deref(p) != 0 {
		t.Error("expected zero value for nil pointer")
	}
}

func TestDerefOr_NonNil(t *testing.T) {
	s := "actual"
	if ptr.DerefOr(&s, "fallback") != "actual" {
		t.Error("expected 'actual'")
	}
}

func TestDerefOr_Nil(t *testing.T) {
	var p *string
	if ptr.DerefOr(p, "fallback") != "fallback" {
		t.Error("expected 'fallback'")
	}
}

func TestOf_Bool(t *testing.T) {
	p := ptr.Of(true)
	if !*p {
		t.Error("expected true")
	}
}
