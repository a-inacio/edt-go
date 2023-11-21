package eventhub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
)

type Handler interface {
	Handler(ctx context.Context, e event.Event) error
}

type callbackHandler struct {
	cb func(ctx context.Context, e event.Event) error
}

type handlers struct {
	callbacks []Handler
}

func (cbh callbackHandler) Handler(ctx context.Context, e event.Event) error {
	return cbh.cb(ctx, e)
}

func ToHandler(cb func(ctx context.Context, e event.Event) error) Handler {
	return &callbackHandler{
		cb: cb,
	}
}
