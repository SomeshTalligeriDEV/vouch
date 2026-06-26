// Package set provides a generic Set type built on a map.
package set

// Set is an unordered collection of unique comparable values.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New creates an empty Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// From creates a Set pre-populated with the given values.
func From[T comparable](values ...T) *Set[T] {
	s := New[T]()
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Add inserts v into the set.
func (s *Set[T]) Add(v T) {
	s.m[v] = struct{}{}
}

// Remove deletes v from the set.
func (s *Set[T]) Remove(v T) {
	delete(s.m, v)
}

// Contains returns true if v is in the set.
func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Slice returns the set elements as an unsorted slice.
func (s *Set[T]) Slice() []T {
	out := make([]T, 0, len(s.m))
	for v := range s.m {
		out = append(out, v)
	}
	return out
}

// Union returns a new set containing all elements from s and other.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := New[T]()
	for v := range s.m {
		result.Add(v)
	}
	for v := range other.m {
		result.Add(v)
	}
	return result
}

// Intersection returns a new set containing only elements present in both sets.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := New[T]()
	for v := range s.m {
		if other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns a new set of elements in s but not in other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := New[T]()
	for v := range s.m {
		if !other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}
