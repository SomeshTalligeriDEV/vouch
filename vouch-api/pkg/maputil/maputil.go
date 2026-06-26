// Package maputil provides generic map helper functions.
package maputil

// Keys returns all keys of m in unspecified order.
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Values returns all values of m in unspecified order.
func Values[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// Merge returns a new map containing all key-value pairs from all maps.
// Later maps' values take precedence on duplicate keys.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	out := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// Filter returns a new map containing only entries for which keep returns true.
func Filter[K comparable, V any](m map[K]V, keep func(K, V) bool) map[K]V {
	out := make(map[K]V)
	for k, v := range m {
		if keep(k, v) {
			out[k] = v
		}
	}
	return out
}

// MapValues returns a new map with the same keys but values transformed by fn.
func MapValues[K comparable, V, R any](m map[K]V, fn func(V) R) map[K]R {
	out := make(map[K]R, len(m))
	for k, v := range m {
		out[k] = fn(v)
	}
	return out
}

// Invert returns a new map with keys and values swapped.
// If values are not unique, later entries overwrite earlier ones.
func Invert[K, V comparable](m map[K]V) map[V]K {
	out := make(map[V]K, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

// GetOrDefault returns the value for key k, or fallback if k is not present.
func GetOrDefault[K comparable, V any](m map[K]V, k K, fallback V) V {
	if v, ok := m[k]; ok {
		return v
	}
	return fallback
}
