// Package search provides lightweight fuzzy and substring search helpers.
package search

import (
	"strings"
	"unicode"
)

// Normalize lowercases and strips diacritics-like noise for comparison.
func Normalize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// Contains reports whether haystack contains needle after normalizing both.
func Contains(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	return strings.Contains(Normalize(haystack), Normalize(needle))
}

// Score returns a relevance score (higher = more relevant) for a query against a text.
// Returns 0 if there is no match.
func Score(text, query string) float64 {
	if query == "" {
		return 1.0
	}
	nText := Normalize(text)
	nQuery := Normalize(query)
	if nText == nQuery {
		return 1.0
	}
	if strings.HasPrefix(nText, nQuery) {
		return 0.9
	}
	if strings.Contains(nText, nQuery) {
		// Score by how early the match appears.
		idx := strings.Index(nText, nQuery)
		position := 1.0 - float64(idx)/float64(len(nText)+1)
		return 0.5 + position*0.3
	}
	// Check if all words in the query appear in the text.
	words := strings.Fields(nQuery)
	allMatch := true
	for _, w := range words {
		if !strings.Contains(nText, w) {
			allMatch = false
			break
		}
	}
	if allMatch {
		return 0.4
	}
	return 0
}

// FilterSlice returns only the strings from items that contain query (case-insensitive).
func FilterSlice(items []string, query string) []string {
	if query == "" {
		return items
	}
	var result []string
	for _, item := range items {
		if Contains(item, query) {
			result = append(result, item)
		}
	}
	return result
}
