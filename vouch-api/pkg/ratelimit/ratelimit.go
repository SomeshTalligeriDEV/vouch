// Package ratelimit provides a simple in-process token-bucket rate limiter.
package ratelimit

import (
	"sync"
	"time"
)

// Bucket is a token-bucket rate limiter for a single key.
type Bucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64 // tokens per second
	lastFill time.Time
}

// NewBucket creates a bucket with the given capacity and refill rate (per second).
func NewBucket(capacity, ratePerSecond float64) *Bucket {
	return &Bucket{
		tokens:   capacity,
		capacity: capacity,
		rate:     ratePerSecond,
		lastFill: time.Now(),
	}
}

// Allow returns true and consumes one token if a request is permitted.
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// AllowN returns true if n tokens are available and consumes them.
func (b *Bucket) AllowN(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= n {
		b.tokens -= n
		return true
	}
	return false
}

func (b *Bucket) refill() {
	now := time.Now()
	elapsed := now.Sub(b.lastFill).Seconds()
	b.tokens = min(b.capacity, b.tokens+elapsed*b.rate)
	b.lastFill = now
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Map is a thread-safe collection of rate-limit buckets keyed by string.
type Map struct {
	mu       sync.Mutex
	buckets  map[string]*Bucket
	capacity float64
	rate     float64
}

// NewMap creates a Map where each key gets a bucket with the given capacity/rate.
func NewMap(capacity, ratePerSecond float64) *Map {
	return &Map{
		buckets:  make(map[string]*Bucket),
		capacity: capacity,
		rate:     ratePerSecond,
	}
}

// Allow returns true if the given key is below the rate limit.
func (m *Map) Allow(key string) bool {
	m.mu.Lock()
	b, ok := m.buckets[key]
	if !ok {
		b = NewBucket(m.capacity, m.rate)
		m.buckets[key] = b
	}
	m.mu.Unlock()
	return b.Allow()
}
