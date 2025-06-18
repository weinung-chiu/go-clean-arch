package adapter

import (
	"context"
	"go-clean-arch/internal/usecase"
	"time"
)

type ChaosEventBroker struct {
	wrapped      usecase.SessionEventBus
	enableEcho   bool
	enableTick   bool
	tickInterval time.Duration

	lastSessionID string // For testing purposes, to track the last session ID
}

func NewChaosEventBroker(wrapped usecase.SessionEventBus) *ChaosEventBroker {
	ceb := &ChaosEventBroker{
		wrapped:      wrapped,
		enableEcho:   true,             // Enable echo by default
		enableTick:   true,             // Enable ticker by default
		tickInterval: 13 * time.Second, // Default ticker interval
	}
	go ceb.startTicker()
	return ceb
}

func (c *ChaosEventBroker) Publish(ctx context.Context, event *usecase.SessionEvent) error {
	c.lastSessionID = event.SessionID // Store the last session ID for testing purposes
	_ = c.wrapped.Publish(ctx, &usecase.SessionEvent{
		Type:      usecase.SessionEventBroadcast,
		SessionID: event.SessionID,
		Payload:   "dummy payload for chaos event",
		Timestamp: time.Now(),
	})
	return c.wrapped.Publish(ctx, event)
}

func (c *ChaosEventBroker) Subscribe(ctx context.Context, sessionID string) (<-chan *usecase.SessionEvent, error) {
	return c.wrapped.Subscribe(ctx, sessionID)
}

func (c *ChaosEventBroker) startTicker() {
	ticker := time.NewTicker(c.tickInterval)
	defer ticker.Stop()
	for {
		<-ticker.C
		event := &usecase.SessionEvent{
			Type:      "SessionEventBroadcast",
			SessionID: c.lastSessionID,
			Payload:   "ticker event from chaos event broker",
			Timestamp: time.Now(),
		}
		// Use background context for broadcast
		_ = c.wrapped.Publish(context.Background(), event)
	}
}
