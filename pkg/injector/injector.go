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

	value, err := i.getValue(t)
	if err != nil {
		return nil, err
	}

	tv := reflect.TypeOf(value)

	if t.Kind() != reflect.Interface && tv.Kind() == reflect.Ptr {
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

func GetValueFromContext[T any](ctx context.Context) (*T, error) {
	return GetValue[T](FromContext(ctx))
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

func (i *Injector) getValues(tt []reflect.Type) ([]interface{}, error) {
	var values []interface{}

	for _, t := range tt {
		value, err := i.getValue(t)

		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func (i *Injector) getValue(t reflect.Type) (interface{}, error) {
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
	fn, t := getFuncReturnValue(value)

	fnArgs := getArgTypes(value)

	i.types = append(i.types, t)

	key := getName(t)

	if len(fnArgs) == 0 {
		i.data[key] = func() interface{} {
			return i.satisfyDependencies(fn, nil)
		}
	} else {
		i.data[key] = func() interface{} {
			values, err := i.getValues(fnArgs)

			if err != nil {
				return err
			}

			return i.satisfyDependencies(fn, values)
		}
	}
}

func getFuncReturnValue(factory interface{}) (reflect.Value, reflect.Type) {
	fn := reflect.ValueOf(factory)
	returnType := fn.Type().Out(0)
	return fn, returnType
}
