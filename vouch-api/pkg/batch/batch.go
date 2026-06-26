// Package batch provides helpers for processing slices in fixed-size chunks.
package batch

// Of splits items into chunks of at most size n.
// The last chunk may be smaller than n.
func Of[T any](items []T, size int) [][]T {
	if size <= 0 {
		size = 1
	}
	if len(items) == 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(items)+size-1)/size)
	for len(items) > 0 {
		end := size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[:end])
		items = items[end:]
	}
	return chunks
}

// Do calls fn for each chunk of items, stopping on first error.
func Do[T any](items []T, size int, fn func(chunk []T) error) error {
	for _, chunk := range Of(items, size) {
		if err := fn(chunk); err != nil {
			return err
		}
	}
	return nil
}

// Map applies fn to each item and returns a slice of results.
// Processing stops on first error.
func Map[T, R any](items []T, size int, fn func(chunk []T) ([]R, error)) ([]R, error) {
	var out []R
	for _, chunk := range Of(items, size) {
		r, err := fn(chunk)
		if err != nil {
			return out, err
		}
		out = append(out, r...)
	}
	return out, nil
}
