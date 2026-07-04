// Package rates implements the read/stream use cases for exchange rates: a
// cached catalogue read path, an in-process SSE hub, and a runner that keeps
// rates fresh (a stand-in for the exchanger parsers).
package rates

import (
	"sync"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
)

// Hub is an in-process fan-out of rate updates, keyed by direction slug. It
// backs the SSE endpoint. A multi-instance deployment would layer Redis pub/sub
// underneath, but the subscriber contract stays the same.
type Hub struct {
	mu   sync.RWMutex
	subs map[string]map[chan []domain.RateRow]struct{}
}

// NewHub builds an empty Hub.
func NewHub() *Hub {
	return &Hub{subs: make(map[string]map[chan []domain.RateRow]struct{})}
}

// Subscribe registers a buffered channel for a direction's updates.
func (h *Hub) Subscribe(dir string) chan []domain.RateRow {
	ch := make(chan []domain.RateRow, 4)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subs[dir] == nil {
		h.subs[dir] = make(map[chan []domain.RateRow]struct{})
	}
	h.subs[dir][ch] = struct{}{}
	return ch
}

// Unsubscribe removes and closes a previously subscribed channel.
func (h *Hub) Unsubscribe(dir string, ch chan []domain.RateRow) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set := h.subs[dir]; set != nil {
		if _, ok := set[ch]; ok {
			delete(set, ch)
			close(ch)
		}
		if len(set) == 0 {
			delete(h.subs, dir)
		}
	}
}

// Publish delivers rows to every subscriber of a direction. Slow subscribers are
// skipped (non-blocking send) rather than stalling the runner.
func (h *Hub) Publish(dir string, rows []domain.RateRow) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs[dir] {
		select {
		case ch <- rows:
		default:
		}
	}
}
