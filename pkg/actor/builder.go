package actor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Builder struct {
	actions   []action.Action
	loopDelay time.Duration
}

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

func (builder *Builder) Do(ctx context.Context) (action.Result, error) {
	return builder.Build().Do(ctx)
}
