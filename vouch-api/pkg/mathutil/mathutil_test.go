package mathutil_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/mathutil"
)

func TestRoundTo(t *testing.T) {
	if mathutil.RoundTo(3.14159, 2) != 3.14 {
		t.Error("expected 3.14")
	}
	if mathutil.RoundTo(2.555, 2) != 2.56 {
		t.Errorf("expected 2.56, got %f", mathutil.RoundTo(2.555, 2))
	}
}

func TestPercent(t *testing.T) {
	if mathutil.Percent(50, 200) != 25.0 {
		t.Errorf("expected 25.0, got %f", mathutil.Percent(50, 200))
	}
	if mathutil.Percent(0, 0) != 0 {
		t.Error("expected 0 for zero total")
	}
}

func TestLerp(t *testing.T) {
	if mathutil.Lerp(0, 100, 0.5) != 50 {
		t.Error("expected 50")
	}
	if mathutil.Lerp(0, 100, 0) != 0 {
		t.Error("expected 0")
	}
	if mathutil.Lerp(0, 100, 1) != 100 {
		t.Error("expected 100")
	}
}

func TestAbs(t *testing.T) {
	if mathutil.Abs(-5) != 5 {
		t.Error("expected 5")
	}
	if mathutil.Abs(5) != 5 {
		t.Error("expected 5")
	}
}

func TestMin(t *testing.T) {
	if mathutil.Min(3, 7) != 3 {
		t.Error("expected 3")
	}
	if mathutil.Min(7, 3) != 3 {
		t.Error("expected 3")
	}
}

func TestMax(t *testing.T) {
	if mathutil.Max(3, 7) != 7 {
		t.Error("expected 7")
	}
}

func TestSum(t *testing.T) {
	if mathutil.Sum(1, 2, 3) != 6 {
		t.Error("expected 6")
	}
	if mathutil.Sum() != 0 {
		t.Error("expected 0 for empty")
	}
}

func TestAverage(t *testing.T) {
	if mathutil.Average(10, 20, 30) != 20 {
		t.Error("expected 20")
	}
	if mathutil.Average() != 0 {
		t.Error("expected 0 for empty")
	}
}
