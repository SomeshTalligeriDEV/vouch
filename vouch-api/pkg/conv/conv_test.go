package conv_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/conv"
)

func TestStringToInt(t *testing.T) {
	if conv.StringToInt("42", 0) != 42 {
		t.Error("expected 42")
	}
	if conv.StringToInt("bad", 7) != 7 {
		t.Error("expected fallback 7")
	}
	if conv.StringToInt("  10 ", 0) != 10 {
		t.Error("expected 10 after trimming whitespace")
	}
}

func TestStringToFloat64(t *testing.T) {
	if conv.StringToFloat64("3.14", 0) != 3.14 {
		t.Error("expected 3.14")
	}
	if conv.StringToFloat64("x", 1.5) != 1.5 {
		t.Error("expected fallback 1.5")
	}
}

func TestStringToBool(t *testing.T) {
	for _, s := range []string{"true", "1", "yes", "TRUE", "YES"} {
		if !conv.StringToBool(s, false) {
			t.Errorf("expected true for %q", s)
		}
	}
	for _, s := range []string{"false", "0", "no"} {
		if conv.StringToBool(s, true) {
			t.Errorf("expected false for %q", s)
		}
	}
	if !conv.StringToBool("unknown", true) {
		t.Error("expected fallback true")
	}
}

func TestIntToString(t *testing.T) {
	if conv.IntToString(99) != "99" {
		t.Error("expected '99'")
	}
}

func TestFloat64ToString(t *testing.T) {
	if conv.Float64ToString(3.14159) != "3.14" {
		t.Errorf("expected '3.14', got %q", conv.Float64ToString(3.14159))
	}
}

func TestClamp(t *testing.T) {
	if conv.Clamp(5, 1, 10) != 5 {
		t.Error("expected 5")
	}
	if conv.Clamp(0, 1, 10) != 1 {
		t.Error("expected clamped to min=1")
	}
	if conv.Clamp(15, 1, 10) != 10 {
		t.Error("expected clamped to max=10")
	}
}

func TestClampFloat64(t *testing.T) {
	if conv.ClampFloat64(0.5, 1.0, 10.0) != 1.0 {
		t.Error("expected clamped to min=1.0")
	}
	if conv.ClampFloat64(15.0, 1.0, 10.0) != 10.0 {
		t.Error("expected clamped to max=10.0")
	}
}

func TestBoolToInt(t *testing.T) {
	if conv.BoolToInt(true) != 1 {
		t.Error("expected 1 for true")
	}
	if conv.BoolToInt(false) != 0 {
		t.Error("expected 0 for false")
	}
}
