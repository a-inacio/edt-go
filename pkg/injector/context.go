package injector

import (
	"context"
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

	i, ok := ctx.Value(reflect.TypeOf(Injector{}).PkgPath()).(*Injector)

	if !ok || i == nil {
		return WithContext(nil)
	}

	return i
}

func GetValueFromContext[T any](ctx context.Context) (*T, error) {
	return Get[T](FromContext(ctx))
}
