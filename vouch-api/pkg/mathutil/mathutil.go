// Package mathutil provides numeric utility functions.
package mathutil

import "math"

// RoundTo rounds f to n decimal places.
func RoundTo(f float64, n int) float64 {
	p := math.Pow10(n)
	return math.Round(f*p) / p
}

// Percent returns what percentage part is of total (0–100).
// Returns 0 if total is zero.
func Percent(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return RoundTo((part/total)*100, 2)
}

// Lerp linearly interpolates between a and b by t (0.0–1.0).
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Abs returns the absolute value of n.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Min returns the smaller of a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of a and b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Sum returns the sum of all values.
func Sum(values ...float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

// Average returns the mean of values, or 0 if empty.
func Average(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Sum(values...) / float64(len(values))
}
