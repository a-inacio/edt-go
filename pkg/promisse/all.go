package promisse

import (
	"context"
	"errors"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/cancellable"
	"reflect"
	"sync"
)

// AllPromise encapsulates a new Promise for a complete fulfillment of one or more given Actions.
// It must be chained back to a Promise by one of the following strategies:
// - Wait: all actions must be fulfilled and completed.
// - WaitWithBailout: all actions must be fulfilled, as soon as one of them fails the Promise fails and execution continues without waiting for the others.
// - WaitWithCancel: all actions must be fulfilled, as soon as one of them fails the Promise fails, attempts to Cancel ongoing actions and execution continues.
type AllPromise struct {
	parent      *Promise
	wg          sync.WaitGroup
	cancellable []*cancellable.Cancellable
	res         []action.Result
	err         []error
}

// All creates a new Promise for the complete fulfillment of one or more given Actions.
func (p *Promise) All(actions ...action.Action) *AllPromise {
	a := &AllPromise{
		parent:      p,
		cancellable: make([]*cancellable.Cancellable, len(actions)),
		res:         make([]action.Result, len(actions)),
		err:         make([]error, len(actions)),
	}

	for i, action := range actions {
		a.cancellable[i] = cancellable.
			NewBuilder().
			FromAction(action).
			WithWaitGroup(&a.wg).
			Build()
	}

	return a
}

// Wait waits for all actions to complete.
func (a *AllPromise) Wait() *Promise {
	then := Future(func(ctx context.Context) (action.Result, error) {
		if a.parent.err != nil {
			return nil, a.parent.err
		}

		completionCh, errCh := a.executeInParallel(ctx)
		errorFound := false

		select {
		case err := <-errCh:
			if err != nil {
				errorFound = true
			}
		case <-completionCh:
			if !errorFound {
				return a.res, nil
			}
		}

		// If we got here, it means that at leas one error was found
		return nil, errors.Join(a.err...)
	})

	then.root = a.parent.root
	a.parent.next = then

	return then
}

// WaitWithBailout waits for all actions to complete, as soon as one of them fails the Promise fails and execution continues without waiting for the others.
func (a *AllPromise) WaitWithBailout() *Promise {
	then := Future(func(ctx context.Context) (action.Result, error) {
		if a.parent.err != nil {
			return nil, a.parent.err
		}

		completionCh, errCh := a.executeInParallel(ctx)

		select {
		case err := <-errCh:
			return nil, err
		case <-completionCh:
			return a.res, nil
		}
	})

	then.root = a.parent.root
	a.parent.next = then

	return then
}

// WaitWithCancel waits for all actions to complete, as soon as one of them fails the Promise fails, attempts to Cancel ongoing actions and execution continues.
func (a *AllPromise) WaitWithCancel() *Promise {
	then := Future(func(ctx context.Context) (action.Result, error) {
		if a.parent.err != nil {
			return nil, a.parent.err
		}

		completionCh, errCh := a.executeInParallel(ctx)

		select {
		case err := <-errCh:
			for _, c := range a.cancellable {
				c.Cancel()
			}

			return nil, err
		case <-completionCh:
			return a.res, nil
		}
	})

	then.root = a.parent.root
	a.parent.next = then

	return then
}

func (a *AllPromise) executeInParallel(ctx context.Context) (completionCh chan any, errCh chan error) {
	chainedCtx := context.WithValue(ctx, reflect.TypeOf(Promise{}).PkgPath(), a.parent.res)

	completionCh = make(chan any, 1)
	errCh = make(chan error, 1)

	for i, c := range a.cancellable {
		// This is important to be made!
		// Otherwise, the closure will capture the last values and not the value of c and i at the time of the iteration
		cb := c
		idx := i

		go func() {
			a.res[idx], a.err[idx] = cb.Do(chainedCtx)
			if a.err[idx] != nil {
				errCh <- a.err[idx]
			}
		}()
	}

	go func() {
		a.wg.Wait()
		close(completionCh)
		close(errCh)
	}()

	return completionCh, errCh
}
