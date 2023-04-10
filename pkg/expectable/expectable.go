package expectable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
)

func NewExpectable(h *event_hub.Hub, e event.Event) *Expectable {
	return &Expectable{h: h, e: e}
}

func (ex *Expectable) Go(ctx context.Context) (action.Result, error) {
	var currentCtx context.Context = nil
	var cancel context.CancelFunc = nil

	if ctx == nil {
		currentCtx = context.Background()
	} else {
		currentCtx = ctx
	}

	h := ex.subscribe()

	defer ex.unsubscribe(h)

	if ex.timeout > 0 {
		// Create a child context that is cancelled when the parent context is cancelled
		currentCtx, cancel = context.WithCancel(currentCtx)
		defer cancel()

		go awaitable.RunAfter(currentCtx, ex.timeout, func(ctx context.Context) (action.Result, error) {
			cancel()
			return action.Nothing()
		})
	}

	for {
		select {
		case <-currentCtx.Done():
			return nil, currentCtx.Err()
		case res := <-h.ch:
			if ex.criteria != nil {
				if ex.criteria(res.Result) {
					return res.Result, res.error
				}
			} else {
				return res.Result, res.error
			}
		}
	}
}

// ==============================================================================
// Auxiliary
// ==============================================================================

func (h *expectableEventHandler) Handler(ctx context.Context, e event.Event) error {
	h.ch <- struct {
		action.Result
		error
	}{
		e,
		nil,
	}
	return nil
}

func (ex *Expectable) subscribe() *expectableEventHandler {
	h := &expectableEventHandler{
		ch: make(chan struct {
			action.Result
			error
		}, 1),
	}

	ex.h.Subscribe(ex.e, h)

	return h
}

func (ex *Expectable) unsubscribe(h *expectableEventHandler) {
	ex.h.Unsubscribe(ex.e, h)
	close(h.ch)
}
