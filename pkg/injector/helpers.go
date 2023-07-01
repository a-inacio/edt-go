package injector

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func fromContext(ctx context.Context) (*Injector, bool) {
	i, ok := ctx.Value(reflect.TypeOf(Injector{}).PkgPath()).(*Injector)
	return i, ok
}

func toContext(parent context.Context, i *Injector) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	return context.WithValue(parent, reflect.TypeOf(Injector{}).PkgPath(), i)
}

func satisfyWithAnotherContext[T any](i *Injector, f interface{}, ctx context.Context) (*T, error) {
	t := reflect.TypeOf(f)

	if isTypeFunc(t) {
		value := i.getSatisfiedInterfaceProxy(f, nil)()

		if value == nil {
			return nil, nil
		}

		return castValue[T](value, t)
	} else {
		return getValueWithContext[T](i, ctx)
	}
}

func getValueWithContext[T any](i *Injector, ctx context.Context) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()

	value, err := i.getValue(t, ctx)
	if err != nil {
		return nil, err
	}

	return castValue[T](value, t)
}

func (i *Injector) satisfyDependencies(fn reflect.Value, args []interface{}) interface{} {
	var paramValues []reflect.Value = nil

	if args != nil {
		numArgs := len(args)
		paramValues = make([]reflect.Value, numArgs)
		for i := 0; i < numArgs; i++ {
			paramValues[i] = reflect.ValueOf(args[i])
		}
	}

	returnValues := fn.Call(paramValues)
	return returnValues[0].Interface()
}

func (i *Injector) getValues(tt []reflect.Type, ctx context.Context) ([]interface{}, error) {
	var values []interface{}

	for _, t := range tt {
		value, err := i.getValue(t, ctx)

		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func (i *Injector) getValue(t reflect.Type, ctx context.Context) (interface{}, error) {
	key := getName(t)

	isInterface := t.Kind() == reflect.Interface

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
		} else if isTypeContext(t) {
			if ctx == nil {
				ctx = i.Context()
			}

			return ctx, nil
		} else {
			return nil, fmt.Errorf("dependency not found %s", key)
		}
	}

	return getter(), nil
}

func getName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		// We don't want the `*` in case this is a pointer  otherwise if passing
		// by a pointer or value the events will be different.
		return strings.TrimPrefix(t.Elem().String(), "*")
	}

	return t.String()
}

func isTypeContext(t reflect.Type) bool {
	return t == reflect.TypeOf((*context.Context)(nil)).Elem()
}

func isTypeFunc(t reflect.Type) bool {
	if t.Kind() != reflect.Func {
		return false
	}
	name := t.String()
	if !strings.HasPrefix(name, "func(") {
		return false
	}

	return true
}

func getArgTypes(f interface{}) []reflect.Type {
	fn := reflect.ValueOf(f)
	numArgs := fn.Type().NumIn()
	argTypes := make([]reflect.Type, numArgs)

	for i := 0; i < numArgs; i++ {
		argTypes[i] = fn.Type().In(i)
	}

	return argTypes
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

func (i *Injector) setSingletonValue(value interface{}, t reflect.Type) {
	i.types = append(i.types, t)

	key := getName(t)

	i.data[key] = func() interface{} {
		return value
	}
}

func (i *Injector) setSingletonFunc(value interface{}) {
	_, t := getFuncReturnValue(value)
	i.data[getName(t)] = i.getSatisfiedInterfaceProxy(value, nil)
}

func (i *Injector) getSatisfiedInterfaceProxy(value interface{}, ctx context.Context) func() interface{} {
	fn, t := getFuncReturnValue(value)

	fnArgs := getArgTypes(value)

	i.types = append(i.types, t)

	var singleton interface{} = nil

	if len(fnArgs) == 0 {
		return func() interface{} {
			if singleton == nil {
				singleton = i.satisfyDependencies(fn, nil)
			}

			return singleton
		}
	} else {
		return func() interface{} {
			if singleton == nil {
				values, err := i.getValues(fnArgs, ctx)

				if err != nil {
					return err
				}

				singleton = i.satisfyDependencies(fn, values)
			}

			return singleton
		}
	}
}

func getFuncReturnValue(factory interface{}) (reflect.Value, reflect.Type) {
	fn := reflect.ValueOf(factory)
	returnType := fn.Type().Out(0)
	return fn, returnType
}

func castValue[T any](value interface{}, t reflect.Type) (*T, error) {
	tv := reflect.TypeOf(value)
	tt := reflect.TypeOf((*T)(nil)).Elem()

	if tt.Kind() != reflect.Interface && t.Kind() != reflect.Interface && tv.Kind() == reflect.Ptr {
		// Cast the value to the desired pointer type.
		typedVal, ok := value.(*T)
		if !ok {
			return nil, fmt.Errorf("value for key %s is not of pointer type %T", getName(t), typedVal)
		}

		return typedVal, nil
	} else {
		// Cast the value to the desired type.
		typedVal, ok := value.(T)
		if !ok {
			return nil, fmt.Errorf("value for key %s is not of type %T", getName(t), typedVal)
		}

		return &typedVal, nil
	}
}
