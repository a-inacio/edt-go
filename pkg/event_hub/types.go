package event_hub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/rosetta-logger-go/pkg/logger"
	"sync"
)

type Hub struct {
	mu            sync.Mutex
	l             logger.Logger
	subscriptions map[string]handlers
}

type HubConfig struct {
	Logger logger.Logger
}

type EventHandler interface {
	Handler(ctx context.Context, e event.Event) error
}

type handlers struct {
	callbacks []EventHandler
}
