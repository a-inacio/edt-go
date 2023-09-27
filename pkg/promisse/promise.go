package promisse

import (
	"context"
	"fmt"
	"github.com/a-inacio/edt-go/pkg/action"
	"reflect"
	"sync"
)

type Promise struct {
	ctx context.Context
	wg  sync.WaitGroup
	r   action.Result
	e   error
}

func Future(ctx context.Context, a action.Action) *Promise {
	if ctx == nil {
		ctx = context.Background()
	}

	await := &Promise{ctx: ctx}

	await.wg.Add(1)

	go func(ctx context.Context, a action.Action) {
		defer await.wg.Done()
		await.r, await.e = a(ctx)
	}(ctx, a)

	return await
}

func GetValue[T any](a *Promise) (*T, error) {
	a.wg.Wait()

	if a.e != nil {
		// Execution failed
		return nil, a.e
	}

	t := reflect.TypeOf((*T)(nil)).Elem()

	// Cast the value to the desired type.
	typedVal, ok := a.r.(T)
	if !ok {
		key := t.String()
		return nil, fmt.Errorf("value for key %s is not of type %T", key, typedVal)
	}

	return &typedVal, nil
}
