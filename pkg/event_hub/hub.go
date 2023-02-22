package event_hub

import (
	"context"
	"github.com/a-inacio/rosetta-logger-go/pkg/logger"
	"github.com/a-inacio/rosetta-logger-go/pkg/rosetta"
	"reflect"
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

func (h *Hub) Subscribe(event interface{}, handler EventHandler) {
	typeName := reflect.TypeOf(event).Name()

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[typeName]
	if !contains {
		subscriptions = handlers{callbacks: make([]EventHandler, 0)}
	}

	subscriptions.callbacks = append(subscriptions.callbacks, handler)
	h.subscriptions[typeName] = subscriptions

	h.mu.Unlock()
}

func (h *Hub) Unsubscribe(event interface{}, handler EventHandler) {
	typeName := reflect.TypeOf(event).Name()

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[typeName]
	if contains && len(subscriptions.callbacks) > 0 {
		callbacks := subscriptions.callbacks

		// remove handler
		for idx, v := range callbacks {
			if v == handler {
				callbacks = append(callbacks[0:idx], callbacks[idx+1:]...)
			}
		}

		subscriptions.callbacks = callbacks
		h.subscriptions[typeName] = subscriptions
	}

	h.mu.Unlock()
}

func (h *Hub) Publish(event interface{}, ctx context.Context) {
	log := h.l
	typeName := reflect.TypeOf(event).Name()

	var callbacks []EventHandler

	h.mu.Lock()

	subscriptions, contains := h.subscriptions[typeName]
	if contains && len(subscriptions.callbacks) > 0 {
		callbacks = subscriptions.callbacks
	}

	h.mu.Unlock()

	if callbacks == nil {
		return
	}

	for _, value := range callbacks {
		callback := value // capture the value for the closure!
		go func() {
			err := callback.Handler(ctx, event)
			if err != nil {
				log.Warn("Event handler failed", "reason", err)
			}
		}()
	}
}
