package injector

import (
	"context"
	"fmt"
	"reflect"
)

// Injector is a dependency injection container.
type Injector struct {
	data  map[string]func() interface{}
	types []reflect.Type
	ctx   context.Context
}

// MustGet will attempt to get the value of a type T from the injector, and panic if it is not found.
// Use this method if you want to enforce that a dependency is present avoiding the error checking at the cost of halting execution.
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

// Get will attempt to get the value of a type T from the injector. If the value is not found, it will return an error.
func Get[T any](i *Injector) (*T, error) {
	return getValueWithContext[T](i, nil)
}

func Satisfy[T any](i *Injector, f interface{}) (*T, error) {
	return satisfyWithAnotherContext[T](i, f, nil)
}

// SetSingleton sets a singleton value that will be returned every time the type is requested.
func (i *Injector) SetSingleton(value interface{}) *Injector {
	t := reflect.TypeOf(value)

	if isTypeFunc(t) {
		i.setSingletonFunc(value)
	} else {
		i.setSingletonValue(value, t)
	}

	return i
}

// SetFactory sets a factory function that will be used to create a new instance of the type.
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

// Context returns the context associated with the injector.
func (i *Injector) Context() context.Context {
	return i.ctx
}
