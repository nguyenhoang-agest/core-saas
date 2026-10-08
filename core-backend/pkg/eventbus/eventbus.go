package eventbus

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Handler represents a typed event handler callback.
type Handler[T any] func(ctx context.Context, event T) error

// Bus is a thread-safe, in-process event bus.
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]any
}

// New creates an initialized in-process Bus.
func New() *Bus {
	return &Bus{
		handlers: make(map[string][]any),
	}
}

// Subscribe registers a typed handler for a given topic.
func Subscribe[T any](b *Bus, topic string, handler Handler[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

// Publish dispatches a typed event synchronously to all registered topic handlers.
func Publish[T any](ctx context.Context, b *Bus, topic string, event T) error {
	b.mu.RLock()
	handlers, exists := b.handlers[topic]
	b.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return nil
	}

	var errList []error
	for _, untyped := range handlers {
		handler, ok := untyped.(Handler[T])
		if !ok {
			continue
		}
		if err := handler(ctx, event); err != nil {
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("errors handling event on topic %q: %w", topic, errors.Join(errList...))
	}

	return nil
}

// Topic provides a type-safe handle for publishing and subscribing to a named topic.
type Topic[T any] struct {
	bus  *Bus
	name string
}

// NewTopic creates a type-safe Topic handle.
func NewTopic[T any](bus *Bus, name string) *Topic[T] {
	return &Topic[T]{bus: bus, name: name}
}

// Subscribe binds a handler to this topic.
func (t *Topic[T]) Subscribe(handler Handler[T]) {
	Subscribe(t.bus, t.name, handler)
}

// Publish sends an event to all subscribers of this topic.
func (t *Topic[T]) Publish(ctx context.Context, event T) error {
	return Publish(ctx, t.bus, t.name, event)
}
