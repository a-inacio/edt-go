package loopable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

func (l *Loopable) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		for _, a := range l.actions {
			a(ctx)

			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
		}

		select {
		case <-time.After(l.delay):
			// Wait for a certain delay
		case <-ctx.Done():
			// The context was cancelled, cancel the delay and return the error
			return action.FromError(ctx.Err())
		}
	}
}

func RunForever(ctx context.Context, delay time.Duration, actions ...action.Action) (action.Result, error) {
	return NewBuilder().
		WithDelay(delay).
		LoopOn(actions...).
		Build().
		Go(ctx)
}
