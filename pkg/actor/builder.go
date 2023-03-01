package actor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) LoopingForever(loopDelay time.Duration, actions ...action.Action) *Builder {
	builder.actions = actions
	builder.loopDelay = loopDelay
	return builder
}

func (builder *Builder) Build() *Actor {
	return &Actor{
		actions:   builder.actions,
		loopDelay: builder.loopDelay,
	}
}

func (builder *Builder) Go(ctx context.Context) (action.Result, error) {
	return builder.Build().Go(ctx)
}
