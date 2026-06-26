// Package idgen provides deterministic and random ID generation helpers.
package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Hex generates a random hex string of n bytes (resulting in 2n hex chars).
func Hex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("idgen: crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// Short returns a 12-hex-character random ID (suitable for display).
func Short() string {
	return Hex(6)
}

// Long returns a 32-hex-character random ID.
func Long() string {
	return Hex(16)
}

// Prefixed returns a prefixed ID like "usr_a1b2c3d4e5f6".
func Prefixed(prefix string) string {
	return prefix + "_" + Hex(8)
}

// TimeSortable returns a time-sortable ID: hex(unix_ms) + random suffix.
// IDs generated later will sort lexicographically after earlier ones.
func TimeSortable() string {
	ms := time.Now().UnixMilli()
	return fmt.Sprintf("%016x%s", ms, Hex(8))
}

// IsHex returns true if s is a valid lowercase hex string of exactly n chars.
func IsHex(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}
