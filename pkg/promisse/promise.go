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
	running bool
}

// Future creates a new Promise from the given Action.
func Future(a action.Action) *Promise {
	p := &Promise{}

	p.root = p

	p.cb = cancellable.
		NewBuilder().
		FromAction(a).
		WithWaitGroup(&p.wg).
		Build()

	p.wg.Add(1)

	return p
}

// Then chains a new Promise from the given Action, into an existent Promise.
func (p *Promise) Then(a action.Action) *Promise {
	then := Future(func(ctx context.Context) (action.Result, error) {
		p.wg.Wait()
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

// FromContext returns the chained value from the given context.
// Beware that this method is only applicable within a chained action.
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
// Resolving the promise is a blocking operation and will wait for all the promises to complete (or any to fail).
// If a promise fails to execute, the actual error is returned.
// Getting the value of a promise is an idempotent operation, it will always return the same value.
func ValueOf[T any](a *Promise) (*T, error) {
	lastChild := a

	for lastChild.next != nil {
		lastChild = lastChild.next
	}
	lastChild.wg.Wait()

	if a.root.err != nil {
		// Execution failed
		return nil, a.root.err
	}

	t := reflect.TypeOf((*T)(nil)).Elem()

	// Cast the value to the desired type.
	typedVal, ok := lastChild.res.(T)
	if !ok {
		key := t.String()
		return nil, fmt.Errorf("the promisse result %s is not of type %T", key, typedVal)
	}

	return &typedVal, nil
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

	if ctx == nil {
		ctx = context.Background()
	}

	p.root.running = true
	p.root.execute(ctx)

	lastChild := p.root

	for lastChild.next != nil {
		lastChild = lastChild.next
	}

	lastChild.wg.Wait()

	if p.root.err != nil {
		// Execution failed
		return nil, p.root.err
	}

	return p.res, p.err
}

func (p *Promise) execute(ctx context.Context) {
	go func() {
		p.res, p.err = p.cb.Do(ctx)

		if p.err != nil {
			p.root.err = p.err
			return
		}

		if p.next != nil {
			p.next.execute(ctx)
		}
	}()
}
