package delayable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Builder struct {
	delay  time.Duration
	action action.Action
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) FromAction(action action.Action) *Builder {
	builder.action = action
	return builder
}

func (builder *Builder) WithDelay(delay time.Duration) *Builder {
	builder.delay = delay
	return builder
}

func (builder *Builder) Build() *Delayable {
	operation := builder.action

	return &Delayable{
		delay:     builder.delay,
		operation: operation,
	}
}

func (builder *Builder) Do(ctx context.Context) (action.Result, error) {
	return builder.Build().Do(ctx)
}
