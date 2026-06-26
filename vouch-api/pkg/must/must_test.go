package must_test

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/must"
)

func TestDo_NilError(t *testing.T) {
	must.Do(nil) // should not panic
}

func TestDo_PanicsOnError(t *testing.T) {
	err := must.Recover(func() {
		must.Do(errors.New("boom"))
	})
	if err == nil {
		t.Error("expected panic to be recovered as error")
	}
}

func TestGet_ReturnsValue(t *testing.T) {
	v := must.Get(42, nil)
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestGet_PanicsOnError(t *testing.T) {
	err := must.Recover(func() {
		must.Get("", errors.New("get failed"))
	})
	if err == nil {
		t.Error("expected panic to be recovered")
	}
}

func TestTruthy_Pass(t *testing.T) {
	must.Truthy(true, "should not panic")
}

func TestTruthy_Panics(t *testing.T) {
	err := must.Recover(func() {
		must.Truthy(false, "assertion message")
	})
	if err == nil {
		t.Error("expected panic")
	}
}

func TestNotNil_Pass(t *testing.T) {
	x := 5
	ptr := must.NotNil(&x, "x")
	if *ptr != 5 {
		t.Errorf("expected 5, got %d", *ptr)
	}
}

func TestNotNil_Panics(t *testing.T) {
	err := must.Recover(func() {
		must.NotNil((*int)(nil), "val")
	})
	if err == nil {
		t.Error("expected panic for nil pointer")
	}
}

func TestRecover_NoPanic(t *testing.T) {
	err := must.Recover(func() {})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
