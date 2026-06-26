// Package email provides email address validation and normalization.
package email

import (
	"errors"
	"regexp"
	"strings"
)

// ErrInvalidEmail is returned when an email address fails validation.
var ErrInvalidEmail = errors.New("invalid email address")

// simple RFC 5322-compatible pattern; not exhaustive but practical.
var emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validate returns ErrInvalidEmail if the address is malformed.
func Validate(addr string) error {
	if !emailRE.MatchString(addr) {
		return ErrInvalidEmail
	}
	return nil
}

// Normalize lowercases and trims whitespace from an email address.
func Normalize(addr string) string {
	return strings.ToLower(strings.TrimSpace(addr))
}

// NormalizeAndValidate normalizes then validates the address.
func NormalizeAndValidate(addr string) (string, error) {
	n := Normalize(addr)
	if err := Validate(n); err != nil {
		return "", err
	}
	return n, nil
}
