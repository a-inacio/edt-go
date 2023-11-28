package cancellable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

// Cancellable is a cancellable action
type Cancellable struct {
	action    action.Action
	wg        *sync.WaitGroup
	completed sync.WaitGroup
	res       action.Result
	err       error
	started   sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// Do the action and return the result or error
// Execution is cancelled if the context is cancelled, if the underlying action fails to implement the cancellation, it will be executed to completion in the background. It can still be made an attempt of graceful shutdown by passing a distinct cancellation context on the Wait method.
func (c *Cancellable) Do(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.started.Done()

	defer c.wg.Done()

	c.completed.Add(1)

	ch := make(chan any, 1)

	go func() {
		defer close(ch)
		defer c.completed.Done()

		c.res, c.err = c.action(c.ctx)
	}()

	select {
	case <-ch:
		return c.res, c.err
	case <-c.ctx.Done():
		c.err = c.ctx.Err()
		return action.FromError(c.err)
	}
}

// Wait for the action to complete or for the context to be cancelled
// It can be utilised to wait for the completion or for graceful shutdown purposes, by allowing a distinct cancellation context to be passed.
func (c *Cancellable) Wait(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ch := make(chan any, 1)

	go func() {
		c.completed.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		return c.res, c.err
	case <-ctx.Done():
		c.err = ctx.Err()
		return action.FromError(c.err)
	}
}

// Cancel the action
func (c *Cancellable) Cancel() {
	c.started.Wait()
	c.cancel()
}
