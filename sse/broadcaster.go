package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Event struct {
	Type string      `json:"type"` // e.g., "queue_update"
	Data interface{} `json:"data"`
}

type subscriber chan Event

type Broadcaster struct {
	lock        sync.RWMutex
	subscribers map[uint]map[subscriber]struct{} // businessID -> set of subscribers
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make(map[uint]map[subscriber]struct{}),
	}
}

func (b *Broadcaster) Subscribe(businessID uint) subscriber {
	ch := make(subscriber, 10)
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.subscribers[businessID]; !ok {
		b.subscribers[businessID] = make(map[subscriber]struct{})
	}
	b.subscribers[businessID][ch] = struct{}{}
	return ch
}

func (b *Broadcaster) Unsubscribe(businessID uint, ch subscriber) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if set, ok := b.subscribers[businessID]; ok {
		delete(set, ch)
		close(ch)
	}
}

func (b *Broadcaster) Publish(businessID uint, ev Event) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	set, ok := b.subscribers[businessID]
	if !ok {
		return
	}
	for ch := range set {
		select {
		case ch <- ev:
		default:
			// subscriber full; drop (to keep non-blocking)
			log.Println("Dropped event for subscriber")
		}
	}
}

// helper to produce JSON string for SSE
func (e Event) JSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// convenience
func (b *Broadcaster) PublishQueueUpdate(businessID uint, data interface{}) {
	ev := Event{Type: "queue_update", Data: data}
	b.Publish(businessID, ev)
}

func (b *Broadcaster) DebugPrintSubscribers() string {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return fmt.Sprintf("subs: %v", len(b.subscribers))
}
