// Package slice provides generic slice helper functions.
package slice

// Filter returns a new slice containing only elements for which keep returns true.
func Filter[T any](s []T, keep func(T) bool) []T {
	out := make([]T, 0, len(s))
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Map applies fn to each element and returns a new slice of results.
func Map[T, R any](s []T, fn func(T) R) []R {
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}

// Contains returns true if any element equals target.
func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// Unique returns a new slice with duplicates removed, preserving order.
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	out := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// Chunk splits s into chunks of at most size n.
func Chunk[T any](s []T, n int) [][]T {
	if n <= 0 {
		n = 1
	}
	if len(s) == 0 {
		return nil
	}
	out := make([][]T, 0, (len(s)+n-1)/n)
	for len(s) > 0 {
		end := n
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[:end])
		s = s[end:]
	}
	return out
}

// Flatten merges a slice of slices into a single slice.
func Flatten[T any](ss [][]T) []T {
	total := 0
	for _, s := range ss {
		total += len(s)
	}
	out := make([]T, 0, total)
	for _, s := range ss {
		out = append(out, s...)
	}
	return out
}

// First returns the first element matching pred and true, or the zero value and false.
func First[T any](s []T, pred func(T) bool) (T, bool) {
	for _, v := range s {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
