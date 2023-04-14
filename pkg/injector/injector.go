package injector

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Injector struct {
	data  map[string]func() interface{}
	types []reflect.Type
	ctx   context.Context
}

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

func GetValue[T any](i *Injector) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	isInterface := t.Kind() == reflect.Interface

	key := getName(t)

	// Fetch the value from the map using the key.
	getter, ok := i.data[key]
	if !ok {
		if isInterface {
			targets := getTypesThatImplement(i.types, t)

			if len(targets) != 1 {
				return nil, fmt.Errorf("dependency not found and/or no single entry implements interface %s", key)
			} else {
				ikey := getName(targets[0])
				getter, ok = i.data[ikey]

				if !ok {
					return nil, fmt.Errorf("unable to satisfy dependency implementing interface %s with type %s", key, ikey)
				}
			}
		} else {
			return nil, fmt.Errorf("dependency not found %s", key)
		}
	}

	value := getter()

	tv := reflect.TypeOf(value)

	if !isInterface && tv.Kind() == reflect.Ptr {
		// Cast the value to the desired pointer type.
		typedVal, ok := value.(*T)
		if !ok {
			return nil, fmt.Errorf("value for key %s is not of pointer type %T", key, typedVal)
		}

		return typedVal, nil
	} else {
		// Cast the value to the desired type.
		typedVal, ok := value.(T)
		if !ok {
			return nil, fmt.Errorf("value for key %s is not of type %T", key, typedVal)
		}

		return &typedVal, nil
	}
}

func getName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		// We don't want the `*` in case this is a pointer  otherwise if passing
		// by a pointer or value the events will be different.
		return strings.TrimPrefix(t.Elem().String(), "*")
	}

	return t.String()
}

func getTypesThatImplement(types []reflect.Type, i reflect.Type) []reflect.Type {
	var implementingTypes []reflect.Type
	for _, t := range types {
		if t.Implements(i) {
			implementingTypes = append(implementingTypes, t)
		}
	}
	return implementingTypes
}

func GetValueFromContext[T any](ctx context.Context) (*T, error) {
	return GetValue[T](FromContext(ctx))
}

func (i *Injector) SetSingleton(value interface{}) *Injector {
	t := reflect.TypeOf(value)

	i.types = append(i.types, t)

	key := getName(t)

	i.data[key] = func() interface{} {
		return value
	}

	return i
}

func (i *Injector) SetFactory(factory interface{}) *Injector {
	fn := reflect.ValueOf(factory)
	returnType := fn.Type().Out(0)

	i.types = append(i.types, returnType)

	key := getName(returnType)

	i.data[key] = func() interface{} {
		returnValues := fn.Call(nil)
		return returnValues[0].Interface()
	}

	return i
}

func (i *Injector) Context() context.Context {
	return i.ctx
}
