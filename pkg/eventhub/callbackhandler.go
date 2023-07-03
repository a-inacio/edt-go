package eventhub

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
)

func (cbh callbackHandler) Handler(ctx context.Context, e event.Event) error {
	return cbh.cb(ctx, e)
}

func ToHandler(cb func(ctx context.Context, e event.Event) error) Handler {
	return callbackHandler{
		cb: cb,
	}
}
