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

// All creates a new Promise of fulfillment one or more Actions.
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

		chainedCtx := context.WithValue(ctx, reflect.TypeOf(Promise{}).PkgPath(), a.parent.res)

		for i, c := range a.cancellable {
			// This is important to be made!
			// otherwise the closure will capture the last values and not the value of c and i at the time of the iteration
			cb := c
			idx := i

			go func() {
				a.res[idx], a.err[idx] = cb.Do(chainedCtx)
			}()
		}

		a.wg.Wait()

		// Check for errors in the actions
		for _, err := range a.err {
			if err != nil {
				// If there are at least one error, return a multi error with all of them
				// the operation bellow transforms a slice of errors into a multi error and ignores nil errors
				return nil, errors.Join(a.err...)
			}
		}

		// Otherwise return the results
		return a.res, nil
	})

	then.root = a.parent.root
	a.parent.next = then

	return then
}
