package eventhub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
)

type Handler interface {
	Handler(ctx context.Context, e event.Event) error
}

type ActionHandler interface {
	Handler
	TargetEvent() event.Event
}

type callbackHandler struct {
	cb func(ctx context.Context, e event.Event) error
	e  event.Event
}

type handlers struct {
	callbacks []Handler
}

func (cbh callbackHandler) Handler(ctx context.Context, e event.Event) error {
	return cbh.cb(ctx, e)
}

func (cbh callbackHandler) TargetEvent() event.Event {
	return cbh.e
}

func ToHandler(e event.Event, cb func(ctx context.Context, e event.Event) error) ActionHandler {
	return &callbackHandler{
		cb: cb,
		e:  e,
	}
}
