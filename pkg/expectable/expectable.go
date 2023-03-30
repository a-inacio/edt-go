package expectable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
	"sync"
)

func NewExpectable(h *event_hub.Hub, e event.Event) *Expectable {
	return &Expectable{h: h, e: e}
}

func (ex *Expectable) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	h := ex.subscribe()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		h.wg.Wait()
	}

	ex.unsubscribe(h)

	return action.Nothing()
}

// ==============================================================================
// Auxiliary
// ==============================================================================

func (h *expectableEventHandler) Handler(ctx context.Context, e event.Event) error {
	h.wg.Done()
	return nil
}

func (ex *Expectable) subscribe() *expectableEventHandler {
	var wg sync.WaitGroup
	wg.Add(1)
	h := &expectableEventHandler{
		wg: wg,
	}

	ex.h.Subscribe(ex.e, h)

	return h
}

func (ex *Expectable) unsubscribe(h *expectableEventHandler) {
	ex.h.Unsubscribe(ex.e, h)
}
