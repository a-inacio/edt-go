package eventhub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/rosetta-logger-go/pkg/logger"
	"github.com/a-inacio/rosetta-logger-go/pkg/rosetta"
	"reflect"
	"sync"
)

type EventHub struct {
	mu            sync.Mutex
	l             logger.Logger
	subscriptions map[string]handlers
}

type Config struct {
	Logger logger.Logger
}

// NewEventHub creates a new EventHub instance
func NewEventHub(config *Config) *EventHub {
	logger := rosetta.NewLogger(logger.NullLoggerType)

	if config != nil {
		if config.Logger != nil {
			logger = config.Logger
		}
	}
	return &EventHub{subscriptions: make(map[string]handlers), l: logger}
}

// RegisterHandler registers a handler for an event
func (h *EventHub) RegisterHandler(e event.Event, handler Handler) {
	eventName := event.GetName(e)

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[eventName]
	if !contains {
		subscriptions = handlers{callbacks: make([]Handler, 0)}
	}

	subscriptions.callbacks = append(subscriptions.callbacks, handler)
	h.subscriptions[eventName] = subscriptions

	h.mu.Unlock()
}

// UnregisterHandler unregister a handler for an event
func (h *EventHub) UnregisterHandler(e event.Event, handler Handler) {
	eventName := event.GetName(e)

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[eventName]
	if contains && len(subscriptions.callbacks) > 0 {
		callbacks := subscriptions.callbacks

		// remove handler
		for idx, v := range callbacks {
			if v == handler {
				callbacks = append(callbacks[0:idx], callbacks[idx+1:]...)
			}
		}

		subscriptions.callbacks = callbacks
		h.subscriptions[eventName] = subscriptions
	}

	h.mu.Unlock()
}

// Publish publishes an event
func (h *EventHub) Publish(e event.Event, ctx context.Context) *sync.WaitGroup {
	log := h.l
	eventName := event.GetName(e)

	var wg sync.WaitGroup

	var callbacks []Handler

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[eventName]
	if contains && len(subscriptions.callbacks) > 0 {
		callbacks = subscriptions.callbacks
	}

	h.mu.Unlock()

	if callbacks == nil {
		return &wg
	}

	wg.Add(len(callbacks))

	for _, value := range callbacks {
		callback := value // capture the value for the closure!
		go func() {
			defer wg.Done()
			err := callback.Handler(ctx, e)
			if err != nil {
				log.Warn("Event handler failed", "reason", err)
			}
		}()
	}

	return &wg
}

// Subscribe subscribes to an event
func (h *EventHub) Subscribe(e event.Event, action action.Action) Handler {
	handler := ToHandler(func(ctx context.Context, e event.Event) error {
		if ctx == nil {
			ctx = context.Background()
		}

		// Create a sub context with the event value
		evCtx := context.
			WithValue(ctx, reflect.TypeOf(e).PkgPath(), e)

		_, err := action(evCtx)
		return err
	})

	h.RegisterHandler(e, handler)

	return handler
}
