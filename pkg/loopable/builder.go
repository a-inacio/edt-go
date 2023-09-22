package loopable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Builder struct {
	actions []action.Action
	delay   time.Duration
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) LoopOn(actions ...action.Action) *Builder {
	builder.actions = actions
	return builder
}

func (builder *Builder) WithDelay(delay time.Duration) *Builder {
	builder.delay = delay
	return builder
}

func (builder *Builder) Build() *Loopable {
	return &Loopable{
		actions: builder.actions,
		delay:   builder.delay,
	}
}

func (builder *Builder) Go(ctx context.Context) (action.Result, error) {
	return builder.Build().Go(ctx)
}
