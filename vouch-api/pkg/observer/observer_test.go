package observer_test

import (
	"sync"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/observer"
)

func TestPublish_CallsAllSubscribers(t *testing.T) {
	bus := observer.New[string]()

	var mu sync.Mutex
	var received []string

	bus.Subscribe(func(s string) {
		mu.Lock()
		received = append(received, s)
		mu.Unlock()
	})
	bus.Subscribe(func(s string) {
		mu.Lock()
		received = append(received, s)
		mu.Unlock()
	})

	bus.Publish("hello")

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Errorf("expected 2 received events, got %d", len(received))
	}
}

func TestUnsubscribe_RemovesHandler(t *testing.T) {
	bus := observer.New[int]()

	count := 0
	unsub := bus.Subscribe(func(n int) { count++ })

	bus.Publish(1)
	if count != 1 {
		t.Fatalf("expected 1 call, got %d", count)
	}

	unsub()
	bus.Publish(2)
	if count != 1 {
		t.Errorf("expected count to stay 1 after unsubscribe, got %d", count)
	}
}

func TestLen_TracksSubscribers(t *testing.T) {
	bus := observer.New[bool]()
	if bus.Len() != 0 {
		t.Error("expected empty bus")
	}

	unsub1 := bus.Subscribe(func(bool) {})
	unsub2 := bus.Subscribe(func(bool) {})

	if bus.Len() != 2 {
		t.Errorf("expected len=2, got %d", bus.Len())
	}

	unsub1()
	if bus.Len() != 1 {
		t.Errorf("expected len=1 after unsubscribe, got %d", bus.Len())
	}

	unsub2()
	if bus.Len() != 0 {
		t.Errorf("expected len=0 after all unsubscribed, got %d", bus.Len())
	}
}

func TestPublish_NoSubscribers_NoPanic(t *testing.T) {
	bus := observer.New[string]()
	bus.Publish("nobody home")
}

func TestPublish_StructEvent(t *testing.T) {
	type ScoreEvent struct {
		UserID string
		Score  float64
	}

	bus := observer.New[ScoreEvent]()
	var got ScoreEvent

	bus.Subscribe(func(e ScoreEvent) { got = e })
	bus.Publish(ScoreEvent{UserID: "u1", Score: 820.5})

	if got.UserID != "u1" || got.Score != 820.5 {
		t.Errorf("unexpected event: %+v", got)
	}
}
