// Package must provides helpers that panic on error — useful in tests and init code.
package must

import "fmt"

// Do panics if err is non-nil.
func Do(err error) {
	if err != nil {
		panic(fmt.Sprintf("must: unexpected error: %v", err))
	}
}

// Get returns v, panicking if err is non-nil.
func Get[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("must: unexpected error: %v", err))
	}
	return v
}

// Truthy panics with msg if condition is false.
func Truthy(condition bool, msg string) {
	if !condition {
		panic("must: assertion failed: " + msg)
	}
}

// NotNil panics if v is nil.
func NotNil[T any](v *T, name string) *T {
	if v == nil {
		panic("must: " + name + " must not be nil")
	}
	return v
}

// Recover wraps fn and converts any panic into an error.
func Recover(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}()
	fn()
	return nil
}
