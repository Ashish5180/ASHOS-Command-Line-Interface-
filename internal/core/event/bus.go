package event

import (
	"reflect"
	"sync"
)

// Handler is a function that processes an event.
type Handler func(event any)

// EventBus is a central hub for publishing and subscribing to events.
// It enables loose coupling between modules (e.g., Task module doesn't need
// to know about the Streak module).
type EventBus struct {
	subscribers map[reflect.Type][]Handler
	mu          sync.RWMutex
}

// NewEventBus creates a new, empty EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[reflect.Type][]Handler),
	}
}

// Subscribe registers a handler for a specific event type.
// eventSample is an instance of the struct you want to listen for (e.g., TaskCompleted{}).
func (eb *EventBus) Subscribe(eventSample any, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	t := reflect.TypeOf(eventSample)
	eb.subscribers[t] = append(eb.subscribers[t], handler)
}

// Publish broadcasts an event to all interested subscribers.
// Each handler is executed in its own goroutine to prevent blocking the caller.
func (eb *EventBus) Publish(event any) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	t := reflect.TypeOf(event)
	if handlers, ok := eb.subscribers[t]; ok {
		for _, handler := range handlers {
			// In a CLI environment, we run handlers synchronously
			// to ensure persistence before the process exits.
			handler(event)
		}
	}
}
