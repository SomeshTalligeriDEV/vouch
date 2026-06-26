// Package money provides currency formatting and conversion utilities.
package money

import (
	"fmt"
	"math"
)

// FormatUSD formats a float64 amount as a USD string.
// e.g. 1234.5 → "$1,234.50"
func FormatUSD(amount float64) string {
	if amount == 0 {
		return "$0.00"
	}
	negative := amount < 0
	if negative {
		amount = -amount
	}

	cents := int64(math.Round(amount * 100))
	dollars := cents / 100
	remainder := cents % 100

	// Add thousands separators.
	dollarsStr := fmt.Sprintf("%d", dollars)
	if len(dollarsStr) > 3 {
		result := make([]byte, 0, len(dollarsStr)+len(dollarsStr)/3)
		for i, c := range dollarsStr {
			if i > 0 && (len(dollarsStr)-i)%3 == 0 {
				result = append(result, ',')
			}
			result = append(result, byte(c))
		}
		dollarsStr = string(result)
	}

	s := fmt.Sprintf("$%s.%02d", dollarsStr, remainder)
	if negative {
		s = "-" + s
	}
	return s
}

// CentsToFloat converts integer cents to a float64 dollar amount.
func CentsToFloat(cents int64) float64 {
	return float64(cents) / 100.0
}

// FloatToCents converts a float64 dollar amount to integer cents.
func FloatToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

// IsValidAmount returns true if the amount is non-negative and has at most 2 decimal places.
func IsValidAmount(amount float64) bool {
	if amount < 0 {
		return false
	}
	cents := amount * 100
	return math.Abs(cents-math.Round(cents)) < 1e-9
}
