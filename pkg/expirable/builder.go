package expirable

import (
	"context"
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

func (builder *Builder) Build() *Expirable {
	return &Expirable{
		timeout:   builder.timeout,
		hooks:     builder.hooks,
		operation: builder.operation,
	}
}

func (builder *Builder) Go(ctx context.Context) (interface{}, error) {
	return builder.Build().Go(ctx)
}
