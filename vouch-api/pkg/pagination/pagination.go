// Package pagination provides helpers for cursor and offset pagination.
package pagination

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Clamp returns the page and limit values clamped to valid ranges.
func Clamp(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return page, limit
}

// Offset returns the zero-based offset for a SQL/MongoDB skip operation.
func Offset(page, limit int) int {
	p, l := Clamp(page, limit)
	return (p - 1) * l
}

// HasNextPage returns true if there are more items after the current page.
func HasNextPage(page, limit, total int) bool {
	_, l := Clamp(page, limit)
	return page*l < total
}
