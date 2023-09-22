package delayable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Delayable struct {
	delay     time.Duration
	operation action.Action
}

func (d *Delayable) Do(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-time.After(d.delay):
		// Wait for a certain duration
	case <-ctx.Done():
		// The context was cancelled
		return action.FromError(ctx.Err())
	}

	return d.operation(ctx)
}

func RunAfter(ctx context.Context, delay time.Duration, a action.Action) (action.Result, error) {
	return NewBuilder().
		FromAction(a).
		WithDelay(delay).Do(ctx)
}

func WaitFor(ctx context.Context, delay time.Duration) (action.Result, error) {
	return NewBuilder().
		FromAction(action.DoNothing).
		WithDelay(delay).Do(ctx)
}
