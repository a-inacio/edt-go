package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"time"
)

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
			return awaitable.RunAfter(ctx, builder.delay, func(ctx context.Context) (action.Result, error) {
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
