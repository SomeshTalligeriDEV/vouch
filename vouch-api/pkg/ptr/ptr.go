// Package ptr provides generic pointer helpers.
package ptr

// Of returns a pointer to the given value.
func Of[T any](v T) *T { return &v }

// Deref returns the value pointed to by p, or the zero value if p is nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOr returns the value pointed to by p, or fallback if p is nil.
func DerefOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}
