package promisse

import (
	"context"
	"fmt"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/cancellable"
	"reflect"
	"sync"
)

// Promise is a promise of a future value.
// It can be chained with other promises.
type Promise struct {
	wg      sync.WaitGroup
	res     action.Result
	err     error
	cb      *cancellable.Cancellable
	root    *Promise
	next    *Promise
	mu      sync.Mutex
	catch   *cancellable.Cancellable
	finally *cancellable.Cancellable
	running bool
}

// Future creates a new Promise from the given Action.
func Future(a action.Action) *Promise {
	p := &Promise{}

	p.root = p

	p.cb = cancellable.
		NewBuilder().
		FromAction(a).
		Build()

	p.wg.Add(1)

	return p
}

// Then chains a new Promise from the given Action, into an existent Promise.
func (p *Promise) Then(a action.Action) *Promise {
	then := Future(func(ctx context.Context) (action.Result, error) {
		if p.err != nil {
			return nil, p.err
		}

		chainedCtx := context.WithValue(ctx, reflect.TypeOf(Promise{}).PkgPath(), p.res)

		return a(chainedCtx)
	})

	then.root = p.root
	p.next = then

	return then
}

// Catch chains an action, that will only be executed, if the previous promise fails.
// If the catch action fails, nothing happens, the error is ignored.
// If the catch action succeeds, the result of the catch action will be propagated to the next promise.
// If this operation is repeated, only the last catch action will be executed.
func (p *Promise) Catch(a action.Action) *Promise {
	p.catch = cancellable.
		NewBuilder().
		FromAction(a).
		Build()

	return p
}

// Finally chains an action, that will always be executed, regardless of any promise outcome.
// If the finally action fails, nothing happens, the error is ignored.
// There can only be one finally action, if this operation is repeated, only the last finally action will be executed.
func (p *Promise) Finally(a action.Action) *Promise {
	p.root.finally = cancellable.
		NewBuilder().
		FromAction(a).
		Build()

	return p
}

// FromContext returns the chained value from the given context.
// Beware that this method is only applicable within a chained action.
func FromContext[T any](ctx context.Context) (*T, error) {
	val := ctx.Value(reflect.TypeOf(Promise{}).PkgPath())

	t := reflect.TypeOf((*T)(nil)).Elem()

	// Cast the value to the desired type.
	typedVal, ok := val.(T)
	if !ok {
		return nil, fmt.Errorf("the chained valueis not of type %s", t.Name())
	}

	return &typedVal, nil
}

// SliceFromContext returns a chained slice of values from the given context.
// Beware that this method is only applicable within a chained action.
func SliceFromContext[T any](ctx context.Context) ([]T, error) {
	res, err := FromContext[action.Result](ctx)

	if err != nil {
		return nil, err
	}

	// Cast the value to slice of action.Result
	sliceOfResults, ok := (*res).([]action.Result)
	if !ok {
		return nil, fmt.Errorf("the promisse result %v is not a slice", res)
	}

	sliceRes := make([]T, len(sliceOfResults))

	for i, r := range sliceOfResults {
		// Cast the value to the desired type.
		typedVal, ok := r.(T)
		if !ok {
			t := reflect.TypeOf((*T)(nil)).Elem()
			key := t.String()
			return nil, fmt.Errorf("the promisse result %v is not a slice of type %s", typedVal, key)
		}
		sliceRes[i] = typedVal
	}

	return sliceRes, nil
}

// ValueOf resolves and returns the value of the promise as the given type, if the promise cannot be converted to the given type an error is returned.
// Resolving the promise is a blocking operation and will wait for all the promises to complete (or any to fail).
// If a promise fails to execute, the actual error is returned.
// Getting the value of a promise is an idempotent operation, it will always return the same value.
func ValueOf[T any](a *Promise) (*T, error) {
	rootPromise := a.root

	rootPromise.wg.Wait()

	if a.root.err != nil {
		// Execution failed
		return nil, a.root.err
	}

	// Cast the value to the desired type.
	typedVal, ok := rootPromise.res.(T)
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem()
		key := t.String()
		return nil, fmt.Errorf("the promisse result %s is not of type %T", key, typedVal)
	}

	return &typedVal, nil
}

func SliceOf[T any](a *Promise) ([]T, error) {
	res, err := ValueOf[action.Result](a)

	if err != nil {
		return nil, err
	}

	// Cast the value to slice of action.Result
	sliceOfResults, ok := (*res).([]action.Result)
	if !ok {
		return nil, fmt.Errorf("the promisse result %v is not a slice", res)
	}

	sliceRes := make([]T, len(sliceOfResults))

	for i, r := range sliceOfResults {
		// Cast the value to the desired type.
		typedVal, ok := r.(T)
		if !ok {
			t := reflect.TypeOf((*T)(nil)).Elem()
			key := t.String()
			return nil, fmt.Errorf("the promisse result %v is not a slice of type %s", typedVal, key)
		}
		sliceRes[i] = typedVal
	}

	return sliceRes, nil

}

// Do is the entry point to execute the promise and return the outcome.
// Execution is cancelled if the context is cancelled
// This operation can only be executed once, if you need to execute it multiple times, use the Future method to create a new promise.
// You can use the ValueOf method to get the value of the promise in an idempotent way or in a deferred manner.
func (p *Promise) Do(ctx context.Context) (action.Result, error) {
	p.root.mu.Lock()
	defer p.root.mu.Unlock()

	if p.root.running {
		return action.FromError(fmt.Errorf("promisse already running or completed"))
	}

	defer p.root.wg.Done()

	if ctx == nil {
		ctx = context.Background()
	}

	p.root.running = true
	p.root.execute(ctx)

	if p.root.finally != nil {
		// Execute finally action
		p.root.finally.Do(ctx)
	}

	if p.root.err != nil {
		// Execution failed
		return nil, p.root.err
	}

	return p.res, p.err
}

func (p *Promise) execute(ctx context.Context) {
	p.res, p.err = p.cb.Do(ctx)

	if p.err != nil {
		if p.catch != nil {
			// error handling in place:
			// - will allow execution to continue
			// - there is a chance to still get a result value
			p.res, _ = p.catch.Do(ctx)

			// clear error
			p.err = nil
		} else {
			// no error handling, halt execution
			p.root.err = p.err
			return
		}
	}

	if p.next != nil {
		p.next.execute(ctx)
	} else {
		p.root.res = p.res
	}
}
