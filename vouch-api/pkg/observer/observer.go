// Package observer implements a simple in-process publish/subscribe event bus.
package observer

import "sync"

// Handler is a function called when an event is published.
type Handler[T any] func(T)

// Bus is a typed publish/subscribe event bus.
type Bus[T any] struct {
	mu       sync.RWMutex
	handlers map[int]Handler[T]
	next     int
}

// New creates an empty Bus.
func New[T any]() *Bus[T] {
	return &Bus[T]{handlers: make(map[int]Handler[T])}
}

// Subscribe registers fn and returns an unsubscribe function.
func (b *Bus[T]) Subscribe(fn Handler[T]) func() {
	b.mu.Lock()
	id := b.next
	b.next++
	b.handlers[id] = fn
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		delete(b.handlers, id)
		b.mu.Unlock()
	}
}

// Publish calls all registered handlers with event in the caller's goroutine.
func (b *Bus[T]) Publish(event T) {
	b.mu.RLock()
	fns := make([]Handler[T], 0, len(b.handlers))
	for _, fn := range b.handlers {
		fns = append(fns, fn)
	}
	b.mu.RUnlock()

	for _, fn := range fns {
		fn(event)
	}
}

// Len returns the number of active subscribers.
func (b *Bus[T]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers)
}
