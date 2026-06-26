package ratelimit_test

import (
	"testing"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/ratelimit"
)

func TestBucket_AllowsUpToCapacity(t *testing.T) {
	b := ratelimit.NewBucket(3, 1) // 3 tokens, 1/sec refill
	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Errorf("expected Allow()=true on request %d", i+1)
		}
	}
	if b.Allow() {
		t.Error("expected Allow()=false when bucket is empty")
	}
}

func TestBucket_RefillsOverTime(t *testing.T) {
	b := ratelimit.NewBucket(1, 10) // 1 token, 10/sec refill
	b.Allow() // drain
	time.Sleep(150 * time.Millisecond)
	if !b.Allow() {
		t.Error("expected bucket to refill after 150ms with 10/sec rate")
	}
}

func TestBucket_AllowN(t *testing.T) {
	b := ratelimit.NewBucket(5, 1)
	if !b.AllowN(5) {
		t.Error("expected AllowN(5) to succeed with 5 tokens")
	}
	if b.AllowN(1) {
		t.Error("expected AllowN(1) to fail when bucket is empty")
	}
}

func TestMap_DifferentKeys(t *testing.T) {
	m := ratelimit.NewMap(1, 1)
	if !m.Allow("alice") {
		t.Error("expected alice to be allowed")
	}
	if m.Allow("alice") {
		t.Error("expected alice to be rate-limited after 1 request")
	}
	// bob should have an independent bucket
	if !m.Allow("bob") {
		t.Error("expected bob to be allowed independently")
	}
}

func TestMap_SameKeySharedBucket(t *testing.T) {
	m := ratelimit.NewMap(2, 1)
	m.Allow("user1")
	m.Allow("user1")
	if m.Allow("user1") {
		t.Error("expected user1 to be rate-limited after 2 requests")
	}
}
