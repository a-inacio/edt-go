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

// Future creates a new Promise from the given Action.
func Future(ctx context.Context, a action.Action) *Promise {
	if ctx == nil {
		ctx = context.Background()
	}

	p := &Promise{ctx: ctx}

	p.wg.Add(1)

	go func(ctx context.Context, a action.Action) {
		defer p.wg.Done()
		p.r, p.e = a(ctx)
	}(ctx, a)

	return p
}

// Then chains a new Promise from the given Action, into an existent Promise.
func (p *Promise) Then(a action.Action) *Promise {
	return Future(p.ctx, func(ctx context.Context) (action.Result, error) {
		p.wg.Wait()
		if p.e != nil {
			return nil, p.e
		}

		chainedCtx := context.WithValue(p.ctx, reflect.TypeOf(Promise{}).PkgPath(), p.r)

		return a(chainedCtx)
	})
}

// FromContext returns the chained value from the given context.
// Use this method within a chained action.
func FromContext[T any](ctx context.Context) (*T, error) {
	val := ctx.Value(reflect.TypeOf(Promise{}).PkgPath())

	t := reflect.TypeOf((*T)(nil)).Elem()

	// Cast the value to the desired type.
	typedVal, ok := val.(T)
	if !ok {
		return nil, fmt.Errorf("the chained valueis not of type %T", t.Name())
	}

	return &typedVal, nil
}

// ValueOf resolves and returns the value of the promise as the given type, if the promise cannot be converted to the given type an error is returned.
func ValueOf[T any](a *Promise) (*T, error) {
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
		return nil, fmt.Errorf("the promisse result %s is not of type %T", key, typedVal)
	}

	return &typedVal, nil
}
