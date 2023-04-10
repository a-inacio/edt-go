package injector

import (
	"context"
	"fmt"
	"reflect"
)

type Instance struct {
	data map[string]func() interface{}
	ctx  context.Context
}

func WithContext(ctx context.Context) *Instance {
	if ctx == nil {
		ctx = context.Background()
	}

	i := &Instance{
		data: make(map[string]func() interface{}),
	}

	i.ctx = context.WithValue(ctx, reflect.TypeOf(Instance{}).PkgPath(), i)

	return i
}

func FromContext(ctx context.Context) *Instance {
	if ctx == nil {
		return WithContext(nil)
	}

	i, ok := ctx.Value(reflect.TypeOf(Instance{}).PkgPath()).(*Instance)

	if !ok || i == nil {
		return WithContext(nil)
	}

	return i
}

func GetValue[T any](i *Instance, value T) (*T, error) {
	key := reflect.TypeOf(value).String()
	// Fetch the value from the map using the key.
	getter, ok := i.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}

	// Cast the value to the desired type.
	typedVal, ok := getter().(T)
	if !ok {
		return nil, fmt.Errorf("value for key %s is not of type %T", key, typedVal)
	}

	return &typedVal, nil
}

func GetValueFromContext[T any](ctx context.Context, value T) (*T, error) {
	return GetValue(FromContext(ctx), value)
}

func (i *Instance) SetSingleton(value interface{}) *Instance {
	key := reflect.TypeOf(value).String()

	i.data[key] = func() interface{} {
		return value
	}

	return i
}

func (i *Instance) SetFactory(factory interface{}) *Instance {
	fn := reflect.ValueOf(factory)
	returnType := fn.Type().Out(0)

	key := returnType.String()

	i.data[key] = func() interface{} {
		returnValues := fn.Call(nil)
		return returnValues[0].Interface()
	}

	return i
}

func (i *Instance) Context() context.Context {
	return i.ctx
}
