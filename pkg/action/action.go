package action

import (
	"context"
	"fmt"
	"reflect"
)

// Result is the result of an Action.
// Anything can be a Result.
type Result interface{}

// Action is a function that returns a Result and an error.
// By receiving a context, it ensures simplicity of usage and opens possibilities for flow control and value propagation.
type Action func(ctx context.Context) (Result, error)

func FromError(err error) (Result, error) {
	return nil, err
}

func FromErrorf(format string, a ...any) (Result, error) {
	return nil, fmt.Errorf(format, a...)
}

func Nothing() (Result, error) {
	return nil, nil
}

func DoNothing(ctx context.Context) (Result, error) {
	return nil, nil
}

// ValueOf converts a Result to the desired type.
// It returns an error if the conversion is not possible.
func ValueOf[T any](result Result) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Cast the value to the desired type.
	typedVal, ok := result.(T)
	if !ok {
		return nil, fmt.Errorf("value is not of type %s", t.String())
	}

	return &typedVal, nil
}

// SliceOf converts a Result to a slice of the desired type.
// It returns an error if the conversion is not possible.
func SliceOf[T any](result Result) ([]*T, error) {
	// Cast the value to slice of action.Result
	sliceOfResults, ok := result.([]Result)
	if !ok {
		return nil, fmt.Errorf("result %v is not a slice", sliceOfResults)
	}

	sliceRes := make([]*T, len(sliceOfResults))

	for i, r := range sliceOfResults {
		// Cast the value to the desired type.
		typedVal, ok := r.(T)
		if !ok {
			t := reflect.TypeOf((*T)(nil)).Elem()
			key := t.String()
			return nil, fmt.Errorf("result %v is not a slice of type %s", typedVal, key)
		}
		sliceRes[i] = &typedVal
	}

	return sliceRes, nil
}
