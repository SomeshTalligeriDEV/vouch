// Package crypto provides cryptographic utilities for Vouch.
package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
)

// RandHex returns n random bytes encoded as a hex string (length = 2n).
func RandHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// MustRandHex returns n random hex bytes, panicking on error.
// Use only in test helpers or one-shot setup.
func MustRandHex(n int) string {
	s, err := RandHex(n)
	if err != nil {
		panic(err)
	}
	return s
}

// RandBase64URL returns n random bytes encoded as a URL-safe base64 string.
func RandBase64URL(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ConstantTimeEqual compares two strings in constant time to prevent timing attacks.
func ConstantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
