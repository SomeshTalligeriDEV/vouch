package diff_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/diff"
)

type Profile struct {
	Name    string
	Email   string
	Website string
	Score   float64
}

func TestStruct_NoChanges(t *testing.T) {
	a := Profile{Name: "Alice", Email: "alice@example.com", Score: 100.0}
	b := Profile{Name: "Alice", Email: "alice@example.com", Score: 100.0}

	changes := diff.Struct(a, b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d: %v", len(changes), changes)
	}
}

func TestStruct_SingleChange(t *testing.T) {
	a := Profile{Name: "Alice", Email: "alice@example.com"}
	b := Profile{Name: "Alice", Email: "new@example.com"}

	changes := diff.Struct(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Name != "Email" {
		t.Errorf("expected Email changed, got %s", changes[0].Name)
	}
	if changes[0].Old != "alice@example.com" {
		t.Errorf("unexpected old value: %v", changes[0].Old)
	}
}

func TestStruct_MultipleChanges(t *testing.T) {
	a := Profile{Name: "Alice", Score: 100.0, Website: "https://a.com"}
	b := Profile{Name: "Bob", Score: 200.0, Website: "https://b.com"}

	changes := diff.Struct(a, b)
	if len(changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(changes))
	}
}

func TestHasChanges_True(t *testing.T) {
	a := Profile{Name: "Alice"}
	b := Profile{Name: "Bob"}
	if !diff.HasChanges(a, b) {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	a := Profile{Name: "Alice"}
	b := Profile{Name: "Alice"}
	if diff.HasChanges(a, b) {
		t.Error("expected HasChanges to return false")
	}
}

func TestFieldNames_ReturnsNames(t *testing.T) {
	a := Profile{Name: "Alice", Email: "a@b.com"}
	b := Profile{Name: "Bob", Email: "a@b.com"}

	names := diff.FieldNames(a, b)
	if len(names) != 1 || names[0] != "Name" {
		t.Errorf("expected [Name], got %v", names)
	}
}

func TestStruct_PointerInputs(t *testing.T) {
	a := &Profile{Name: "Alice"}
	b := &Profile{Name: "Bob"}

	changes := diff.Struct(a, b)
	if len(changes) != 1 {
		t.Errorf("expected 1 change for pointer inputs, got %d", len(changes))
	}
}
