// internal/core/events.go
package core

import (
	"sync"

	"github.com/yourname/claude-task-manager/pkg/models"
)

// EventHandler is a callback for events
type EventHandler func(event *models.Event)

// EventBus manages event distribution
type EventBus struct {
	mu       sync.RWMutex
	handlers map[models.EventType][]EventHandler
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[models.EventType][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type
func (b *EventBus) Subscribe(eventType models.EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish sends an event to all registered handlers
func (b *EventBus) Publish(event *models.Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		go h(event) // async delivery
	}
}
