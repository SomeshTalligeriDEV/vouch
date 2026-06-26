package set_test

import (
	"sort"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/set"
)

func TestNew_Empty(t *testing.T) {
	s := set.New[string]()
	if s.Len() != 0 {
		t.Errorf("expected empty set, got len=%d", s.Len())
	}
}

func TestAdd_Contains(t *testing.T) {
	s := set.New[string]()
	s.Add("a")
	s.Add("b")
	if !s.Contains("a") {
		t.Error("expected set to contain 'a'")
	}
	if s.Contains("c") {
		t.Error("expected set NOT to contain 'c'")
	}
	if s.Len() != 2 {
		t.Errorf("expected len=2, got %d", s.Len())
	}
}

func TestAdd_Deduplicates(t *testing.T) {
	s := set.From("x", "x", "x")
	if s.Len() != 1 {
		t.Errorf("expected len=1 after duplicate adds, got %d", s.Len())
	}
}

func TestRemove(t *testing.T) {
	s := set.From("a", "b")
	s.Remove("a")
	if s.Contains("a") {
		t.Error("expected 'a' to be removed")
	}
	if s.Len() != 1 {
		t.Errorf("expected len=1 after remove, got %d", s.Len())
	}
}

func TestSlice(t *testing.T) {
	s := set.From("c", "a", "b")
	sl := s.Slice()
	sort.Strings(sl)
	if len(sl) != 3 || sl[0] != "a" || sl[1] != "b" || sl[2] != "c" {
		t.Errorf("unexpected slice: %v", sl)
	}
}

func TestUnion(t *testing.T) {
	a := set.From(1, 2, 3)
	b := set.From(3, 4, 5)
	u := a.Union(b)
	if u.Len() != 5 {
		t.Errorf("expected union len=5, got %d", u.Len())
	}
}

func TestIntersection(t *testing.T) {
	a := set.From(1, 2, 3)
	b := set.From(2, 3, 4)
	i := a.Intersection(b)
	if i.Len() != 2 {
		t.Errorf("expected intersection len=2, got %d", i.Len())
	}
	if !i.Contains(2) || !i.Contains(3) {
		t.Error("expected intersection to contain 2 and 3")
	}
}

func TestDifference(t *testing.T) {
	a := set.From("a", "b", "c")
	b := set.From("b", "c", "d")
	d := a.Difference(b)
	if d.Len() != 1 || !d.Contains("a") {
		t.Errorf("expected difference={a}, got %v", d.Slice())
	}
}
