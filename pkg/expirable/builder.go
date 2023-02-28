package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"time"
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) FromOperation(operation func(ctx context.Context) (interface{}, error)) *Builder {
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
		operation = func(ctx context.Context) (interface{}, error) {
			return awaitable.RunAfter(ctx, builder.delay, func() (any, error) {
				return builder.operation(ctx)
			})
		}
	}

	return &Expirable{
		timeout:   builder.timeout,
		hooks:     builder.hooks,
		operation: operation,
	}
}

func (builder *Builder) Go(ctx context.Context) (interface{}, error) {
	return builder.Build().Go(ctx)
}
