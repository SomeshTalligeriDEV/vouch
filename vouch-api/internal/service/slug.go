package service

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

// slugify converts a title into a URL-safe slug with a short random suffix to
// guarantee uniqueness without a round-trip to the database.
func slugify(title string) string {
	base := strings.ToLower(strings.TrimSpace(title))
	base = nonSlugChars.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "item"
	}
	if len(base) > 60 {
		base = strings.Trim(base[:60], "-")
	}
	return base + "-" + randHex(3)
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "000000"[:n*2]
	}
	return hex.EncodeToString(b)
}
