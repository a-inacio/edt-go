package expectable

import (
	"context"
	"errors"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/eventhub"
	"time"
)

type Builder struct {
	e        event.Event
	h        *eventhub.EventHub
	timeout  time.Duration
	criteria func(e event.Event) bool
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) On(h *eventhub.EventHub) *Builder {
	builder.h = h
	return builder
}

func (builder *Builder) Expect(e event.Event) *Builder {
	builder.e = e
	return builder
}

func (builder *Builder) WithTimeout(timeout time.Duration) *Builder {
	builder.timeout = timeout
	return builder
}

func (builder *Builder) Where(criteria func(e event.Event) bool) *Builder {
	builder.criteria = criteria
	return builder
}

func (builder *Builder) Build() (*Expectable, error) {
	if builder.h == nil {
		return nil, errors.New("missing event hub definition")
	}

	if builder.e == nil {
		return nil, errors.New("missing event definition")
	}

	instance := NewExpectable(builder.h, builder.e)
	instance.timeout = builder.timeout
	instance.criteria = builder.criteria

	return instance, nil
}

func (builder *Builder) Go(ctx context.Context) (action.Result, error) {
	instance, err := builder.Build()

	if err != nil {
		return nil, err
	}

	return instance.Go(ctx)
}
