package cursor_test

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/cursor"
)

func TestEncodeDecode_RoundTrip(t *testing.T) {
	original := cursor.Cursor{ID: "abc123", Value: "2024-01-01"}
	encoded, err := cursor.Encode(original)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if encoded == "" {
		t.Fatal("expected non-empty encoded string")
	}

	decoded, err := cursor.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.ID, original.ID)
	}
}

func TestDecode_InvalidBase64(t *testing.T) {
	_, err := cursor.Decode("!!!notvalid###")
	if !errors.Is(err, cursor.ErrInvalidCursor) {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestDecode_EmptyID(t *testing.T) {
	// Encode a cursor without an ID.
	enc, _ := cursor.Encode(cursor.Cursor{ID: ""})
	_, err := cursor.Decode(enc)
	if !errors.Is(err, cursor.ErrInvalidCursor) {
		t.Errorf("expected ErrInvalidCursor for empty ID, got %v", err)
	}
}

func TestDecode_EmptyString(t *testing.T) {
	_, err := cursor.Decode("")
	if !errors.Is(err, cursor.ErrInvalidCursor) {
		t.Errorf("expected ErrInvalidCursor for empty input, got %v", err)
	}
}

func TestEncode_URLSafe(t *testing.T) {
	c := cursor.Cursor{ID: "some-id-12345"}
	enc, err := cursor.Encode(c)
	if err != nil {
		t.Fatal(err)
	}
	// URL-safe base64 should not contain + or /
	for _, ch := range enc {
		if ch == '+' || ch == '/' {
			t.Errorf("encoded cursor contains non-URL-safe character %q", ch)
		}
	}
}
