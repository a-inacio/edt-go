package expectable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
)

func NewExpectable(h *event_hub.Hub, e event.Event) *Expectable {
	return &Expectable{h: h, e: e}
}

func (ex *Expectable) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	h := ex.subscribe()

	defer ex.unsubscribe(h)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-h.ch:
	}

	return action.Nothing()
}

// ==============================================================================
// Auxiliary
// ==============================================================================

func (h *expectableEventHandler) Handler(ctx context.Context, e event.Event) error {
	h.ch <- struct{}{}
	return nil
}

func (ex *Expectable) subscribe() *expectableEventHandler {
	h := &expectableEventHandler{
		ch: make(chan struct {
		}, 1),
	}

	ex.h.Subscribe(ex.e, h)

	return h
}

func (ex *Expectable) unsubscribe(h *expectableEventHandler) {
	ex.h.Unsubscribe(ex.e, h)
	close(h.ch)
}
