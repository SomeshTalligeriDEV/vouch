// Package cursor provides opaque cursor encoding for API pagination.
// Cursors encode a MongoDB ObjectID + sort field so clients can page through
// results without knowing internal sort keys.
package cursor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// ErrInvalidCursor is returned when a cursor string cannot be decoded.
var ErrInvalidCursor = errors.New("invalid cursor")

// Cursor holds the data encoded inside an opaque page cursor.
type Cursor struct {
	ID    string `json:"id"`
	Value any    `json:"v,omitempty"`
}

// Encode serialises a cursor to a URL-safe base64 string.
func Encode(c Cursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("cursor.Encode: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// MustEncode encodes a cursor, panicking on error. For tests / seed scripts.
func MustEncode(c Cursor) string {
	s, err := Encode(c)
	if err != nil {
		panic(err)
	}
	return s
}

// Decode parses an opaque cursor string.
func Decode(s string) (Cursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Cursor{}, ErrInvalidCursor
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return Cursor{}, ErrInvalidCursor
	}
	if c.ID == "" {
		return Cursor{}, ErrInvalidCursor
	}
	return c, nil
}
