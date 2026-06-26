package money_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/money"
)

func TestFormatUSD_Zero(t *testing.T) {
	if got := money.FormatUSD(0); got != "$0.00" {
		t.Errorf("expected '$0.00', got %q", got)
	}
}

func TestFormatUSD_Positive(t *testing.T) {
	cases := []struct {
		amount float64
		want   string
	}{
		{1, "$1.00"},
		{9.99, "$9.99"},
		{100.50, "$100.50"},
		{1234.56, "$1,234.56"},
		{1000000, "$1,000,000.00"},
	}
	for _, tc := range cases {
		got := money.FormatUSD(tc.amount)
		if got != tc.want {
			t.Errorf("FormatUSD(%.2f) = %q, want %q", tc.amount, got, tc.want)
		}
	}
}

func TestFormatUSD_Negative(t *testing.T) {
	got := money.FormatUSD(-5.00)
	if got != "-$5.00" {
		t.Errorf("expected '-$5.00', got %q", got)
	}
}

func TestCentsToFloat(t *testing.T) {
	got := money.CentsToFloat(1050)
	if got != 10.50 {
		t.Errorf("expected 10.50, got %f", got)
	}
}

func TestFloatToCents(t *testing.T) {
	got := money.FloatToCents(10.50)
	if got != 1050 {
		t.Errorf("expected 1050, got %d", got)
	}
}

func TestIsValidAmount_Valid(t *testing.T) {
	valid := []float64{0, 1, 9.99, 100.50, 1234.56}
	for _, v := range valid {
		if !money.IsValidAmount(v) {
			t.Errorf("IsValidAmount(%f): expected true", v)
		}
	}
}

func TestIsValidAmount_Negative(t *testing.T) {
	if money.IsValidAmount(-1) {
		t.Error("expected false for negative amount")
	}
}
