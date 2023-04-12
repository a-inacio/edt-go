package action

import (
	"context"
	"fmt"
	"reflect"
)

func FromError(err error) (Result, error) {
	return nil, err
}

func Nothing() (Result, error) {
	return nil, nil
}

func DoNothing(ctx context.Context) (Result, error) {
	return nil, nil
}

func GetValue[T any](result Result) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Cast the value to the desired type.
	typedVal, ok := result.(T)
	if !ok {
		return nil, fmt.Errorf("value is not of type %s", t.String())
	}

	return &typedVal, nil
}
