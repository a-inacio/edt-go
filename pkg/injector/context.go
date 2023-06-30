package injector

import (
	"context"
	"fmt"
	"reflect"
)

func WithContext(ctx context.Context) *Injector {
	if ctx == nil {
		ctx = context.Background()
	}

	i := &Injector{
		data: make(map[string]func() interface{}),
	}

	i.ctx = context.WithValue(ctx, reflect.TypeOf(Injector{}).PkgPath(), i)

	return i
}

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

func GetFromContext[T any](ctx context.Context) (*T, error) {
	return Get[T](FromContext(ctx))
}

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

func MustSatisfyFromContext[T any](ctx context.Context, f interface{}) T {
	if ctx == nil {
		panic("context cannot be nil")
	}

	i, ok := fromContext(ctx)

	if !ok || i == nil {
		panic("no injector in context")
	}

	return MustSatisfy[T](i, f)
}
