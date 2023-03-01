package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

func (e *Expirable) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a child context that is cancelled when the parent context is cancelled
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a channel that is used to signal completion of the long-running action
	ch := make(chan struct {
		any
		error
	}, 1)

	// Start the long-running action in a separate goroutine
	go func() {
		res, err := e.action(ctx)
		ch <- struct {
			any
			error
		}{res, err}
	}()

	// Wait for the action to complete or for the timeout to expire
	select {
	case <-ctx.Done():
		// The parent context was cancelled, cancel the child context
		return nil, ctx.Err()
	case res := <-ch:
		if res.error != nil {
			e.onErrorCb(ctx, res.error)
			return nil, res.error
		} else {
			e.onSuccessCb(ctx, res.any)
			return res.any, nil
		}
	case <-time.After(e.timeout):
		// The action timed out...
		e.onExpiredCb(ctx)
		// ...cancel the child context
		cancel()
		return nil, context.DeadlineExceeded
	}
}

func (e *Expirable) onCompletedCb(ctx context.Context, res interface{}, err error) {
	if e.hooks.OnCompleted != nil {
		e.hooks.OnCompleted(ctx, res, err)
	}
}

func (e *Expirable) onSuccessCb(ctx context.Context, res interface{}) {
	if e.hooks.OnSuccess != nil {
		e.hooks.OnSuccess(ctx, res)
	}
}

func (e *Expirable) onErrorCb(ctx context.Context, err error) {
	if e.hooks.OnError != nil {
		e.hooks.OnError(ctx, err)
	}
}

func (e *Expirable) onCanceledCb(ctx context.Context) {
	if e.hooks.OnCanceled != nil {
		e.hooks.OnCanceled(ctx)
	}
}

func (e *Expirable) onExpiredCb(ctx context.Context) {
	if e.hooks.OnExpired != nil {
		e.hooks.OnExpired(ctx)
	}
}
