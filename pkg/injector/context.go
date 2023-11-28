package injector

import (
	"context"
	"fmt"
)

// WithContext creates a new injector with the given context.
func WithContext(ctx context.Context) *Injector {
	i := &Injector{
		data: make(map[string]func() interface{}),
	}

	i.ctx = toContext(ctx, i)

	return i
}

// FromContext gets the injector from the given context, or creates a new one if it does not exist.
func FromContext(ctx context.Context) *Injector {
	if ctx == nil {
		return WithContext(nil)
	}

	i, ok := fromContext(ctx)

	if !ok || i == nil {
		return WithContext(nil)
	}

	return i
}

// GetFromContext gets the value from the given context. If the value does not exist, it will return an error.
// Use this method as quick way to get a single value from the context and to avoid having to retrieve the injector.
// If you need to get multiple values from the context, it is recommended to retrieve the injector and use the Get method.
func GetFromContext[T any](ctx context.Context) (*T, error) {
	return Get[T](FromContext(ctx))
}

// MustGetFromContext gets the value from the given context. If the value does not exist, it will panic.
// Use this method as quick way to get a single value from the context and to avoid having to retrieve the injector.
// If you need to get multiple values from the context, it is recommended to retrieve the injector and use the MustGet method.
func MustGetFromContext[T any](ctx context.Context) T {
	if ctx == nil {
		panic("context cannot be nil")
	}

	i, ok := fromContext(ctx)

	if !ok || i == nil {
		panic("no injector in context")
	}

	value, err := Get[T](i)
	if err != nil {
		panic(fmt.Sprintf("missing dependency - %s", err.Error()))
	}
	return *value
}

// MustResolveFromContext will attempt to resolve the arguments of a function with a return a type T, from the given context.
// If all arguments can be satisfied, the function is called, it will panic otherwise.
func MustResolveFromContext[T any](ctx context.Context, f interface{}) T {
	if ctx == nil {
		panic("context cannot be nil")
	}

	i, ok := fromContext(ctx)

	if !ok || i == nil {
		panic("no injector in context")
	}

	return MustResolve[T](i, f)
}
