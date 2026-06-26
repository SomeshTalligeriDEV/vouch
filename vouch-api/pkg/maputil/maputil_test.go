package maputil_test

import (
	"sort"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/maputil"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := maputil.Keys(m)
	sort.Strings(keys)
	if len(keys) != 3 || keys[0] != "a" || keys[2] != "c" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{"x": 10}
	vals := maputil.Values(m)
	if len(vals) != 1 || vals[0] != 10 {
		t.Errorf("unexpected values: %v", vals)
	}
}

func TestMerge_NoOverlap(t *testing.T) {
	a := map[string]int{"a": 1}
	b := map[string]int{"b": 2}
	merged := maputil.Merge(a, b)
	if merged["a"] != 1 || merged["b"] != 2 {
		t.Errorf("unexpected merge result: %v", merged)
	}
}

func TestMerge_LaterWins(t *testing.T) {
	a := map[string]int{"k": 1}
	b := map[string]int{"k": 99}
	merged := maputil.Merge(a, b)
	if merged["k"] != 99 {
		t.Errorf("expected later map to win, got %d", merged["k"])
	}
}

func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := maputil.Filter(m, func(_ string, v int) bool { return v > 1 })
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(got), got)
	}
	if _, ok := got["a"]; ok {
		t.Error("expected 'a' to be filtered out")
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"x": 2, "y": 3}
	got := maputil.MapValues(m, func(v int) int { return v * 2 })
	if got["x"] != 4 || got["y"] != 6 {
		t.Errorf("unexpected MapValues result: %v", got)
	}
}

func TestInvert(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	inv := maputil.Invert(m)
	if inv[1] != "a" || inv[2] != "b" {
		t.Errorf("unexpected inverted map: %v", inv)
	}
}

func TestGetOrDefault_Exists(t *testing.T) {
	m := map[string]int{"key": 42}
	if maputil.GetOrDefault(m, "key", 0) != 42 {
		t.Error("expected 42")
	}
}

func TestGetOrDefault_Missing(t *testing.T) {
	m := map[string]int{}
	if maputil.GetOrDefault(m, "missing", 7) != 7 {
		t.Error("expected fallback 7")
	}
}
