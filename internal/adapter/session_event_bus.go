package adapter

import (
	"context"
	"go-clean-arch/internal/usecase"
	"sync"
)

type MemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan *usecase.EventQuestionUpdated
}

func NewInMemorySessionEventBroker() *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make(map[string][]chan *usecase.EventQuestionUpdated),
	}
}

// Broadcast sends the event to all subscribers of the given session ID.
func (b *MemoryEventBus) Broadcast(ctx context.Context, event *usecase.EventQuestionUpdated) error {
	b.mu.RLock()
	chs := b.subscribers[event.SessionID]
	b.mu.RUnlock()
	for _, ch := range chs {
		select {
		case ch <- event:
		case <-ctx.Done():
			return ctx.Err()
		default:
			// drop if buffer full
		}
	}
	return nil
}
func (b *MemoryEventBus) Subscribe(ctx context.Context, sessionID string) (<-chan *usecase.EventQuestionUpdated, error) {
	ch := make(chan *usecase.EventQuestionUpdated, 16)
	b.mu.Lock()
	b.subscribers[sessionID] = append(b.subscribers[sessionID], ch)
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.mu.Lock()
		subs := b.subscribers[sessionID]
		for i, c := range subs {
			if c == ch {
				// Remove the channel from the slice
				b.subscribers[sessionID] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		b.mu.Unlock()
		close(ch)
	}()

	return ch, nil
}
