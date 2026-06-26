package slice_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/slice"
)

func TestFilter(t *testing.T) {
	got := slice.Filter([]int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 == 0 })
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Errorf("unexpected filter result: %v", got)
	}
}

func TestMap(t *testing.T) {
	got := slice.Map([]int{1, 2, 3}, func(v int) string {
		if v == 1 {
			return "one"
		}
		return "other"
	})
	if len(got) != 3 || got[0] != "one" {
		t.Errorf("unexpected map result: %v", got)
	}
}

func TestContains(t *testing.T) {
	if !slice.Contains([]string{"a", "b", "c"}, "b") {
		t.Error("expected Contains to return true")
	}
	if slice.Contains([]string{"a", "b", "c"}, "z") {
		t.Error("expected Contains to return false")
	}
}

func TestUnique(t *testing.T) {
	got := slice.Unique([]int{1, 2, 2, 3, 1, 3})
	if len(got) != 3 {
		t.Errorf("expected 3 unique items, got %d: %v", len(got), got)
	}
	if got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestChunk(t *testing.T) {
	got := slice.Chunk([]int{1, 2, 3, 4, 5}, 2)
	if len(got) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(got))
	}
	if len(got[2]) != 1 {
		t.Error("expected last chunk to have 1 item")
	}
}

func TestFlatten(t *testing.T) {
	got := slice.Flatten([][]int{{1, 2}, {3}, {4, 5}})
	if len(got) != 5 || got[0] != 1 || got[4] != 5 {
		t.Errorf("unexpected flatten result: %v", got)
	}
}

func TestFirst_Found(t *testing.T) {
	v, ok := slice.First([]int{1, 3, 5, 6, 7}, func(n int) bool { return n%2 == 0 })
	if !ok || v != 6 {
		t.Errorf("expected (6, true), got (%d, %v)", v, ok)
	}
}

func TestFirst_NotFound(t *testing.T) {
	_, ok := slice.First([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 })
	if ok {
		t.Error("expected false when no match")
	}
}
