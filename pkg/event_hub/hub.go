package event_hub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/rosetta-logger-go/pkg/logger"
	"github.com/a-inacio/rosetta-logger-go/pkg/rosetta"
	"sync"
)

func NewHub(config *HubConfig) *Hub {
	logger := rosetta.NewLogger(logger.NullLoggerType)

	if config != nil {
		if config.Logger != nil {
			logger = config.Logger
		}
	}
	return &Hub{subscriptions: make(map[string]handlers), l: logger}
}

func (h *Hub) Subscribe(e event.Event, handler EventHandler) {
	eventName := event.GetName(e)

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[eventName]
	if !contains {
		subscriptions = handlers{callbacks: make([]EventHandler, 0)}
	}

	subscriptions.callbacks = append(subscriptions.callbacks, handler)
	h.subscriptions[eventName] = subscriptions

	h.mu.Unlock()
}

func (h *Hub) Unsubscribe(e event.Event, handler EventHandler) {
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

func (h *Hub) Publish(e event.Event, ctx context.Context) *sync.WaitGroup {
	log := h.l
	eventName := event.GetName(e)

	var wg sync.WaitGroup

	var callbacks []EventHandler

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
