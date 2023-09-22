package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/delayable"
	"time"
)

type Hooks struct {
	Action      func(ctx context.Context) (interface{}, error)
	OnExpired   func(ctx context.Context)
	OnSuccess   func(ctx context.Context, result interface{})
	OnError     func(ctx context.Context, e error)
	OnCanceled  func(ctx context.Context)
	OnCompleted func(ctx context.Context, result interface{}, err error)
}

type Builder struct {
	timeout   time.Duration
	delay     time.Duration
	operation action.Action
	hooks     Hooks
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) FromOperation(operation action.Action) *Builder {
	builder.operation = operation
	return builder
}

func (builder *Builder) WithTimeout(timeout time.Duration) *Builder {
	builder.timeout = timeout
	return builder
}

func (builder *Builder) WithDelay(delay time.Duration) *Builder {
	builder.delay = delay
	return builder
}

func (builder *Builder) Build() *Expirable {
	operation := builder.operation

	if builder.delay > 0 {
		operation = func(ctx context.Context) (action.Result, error) {
			return delayable.RunAfter(ctx, builder.delay, func(ctx context.Context) (action.Result, error) {
				return builder.operation(ctx)
			})
		}
	}

	return &Expirable{
		timeout: builder.timeout,
		hooks:   builder.hooks,
		action:  operation,
	}
}

func (builder *Builder) Go(ctx context.Context) (action.Result, error) {
	return builder.Build().Go(ctx)
}
