package injector

import (
	"context"
	"fmt"
	"reflect"
)

type Injector struct {
	data  map[string]func() interface{}
	types []reflect.Type
	ctx   context.Context
}

func MustGet[T any](i *Injector) T {
	value, err := Get[T](i)
	if err != nil {
		panic(fmt.Sprintf("missing dependency - %s", err.Error()))
	}
	return *value
}

func MustSatisfy[T any](i *Injector, f interface{}) T {
	res, err := satisfyWithAnotherContext[T](i, f, nil)

	if err != nil {
		panic(fmt.Sprintf("unable to satisfy interface - %s", err.Error()))
	}

	return *res
}

func Get[T any](i *Injector) (*T, error) {
	return getValueWithContext[T](i, nil)
}

func Satisfy[T any](i *Injector, f interface{}) (*T, error) {
	return satisfyWithAnotherContext[T](i, f, nil)
}

func (i *Injector) SetSingleton(value interface{}) *Injector {
	t := reflect.TypeOf(value)

	if isTypeFunc(t) {
		i.setSingletonFunc(value)
	} else {
		i.setSingletonValue(value, t)
	}

	return i
}

func (i *Injector) SetFactory(factory interface{}) *Injector {
	if isTypeFunc(reflect.TypeOf(factory)) {
		fn, returnType := getFuncReturnValue(factory)

		i.types = append(i.types, returnType)

		key := getName(returnType)

		i.data[key] = func() interface{} {
			return i.satisfyDependencies(fn, nil)
		}
	}

	return i
}

func (i *Injector) Context() context.Context {
	return i.ctx
}
