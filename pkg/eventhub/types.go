package eventhub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/rosetta-logger-go/pkg/logger"
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

type Handler interface {
	Handler(ctx context.Context, e event.Event) error
}

type handlers struct {
	callbacks []Handler
}

type callbackHandler struct {
	cb func(ctx context.Context, e event.Event) error
}