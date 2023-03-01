package loopable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

func RunForever(ctx context.Context, delay time.Duration, actions ...action.Action) (action.Result, error) {
	for {
		for _, a := range actions {
			a(ctx)

			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
		}

		select {
		case <-time.After(delay):
			// Wait for a certain delay
		case <-ctx.Done():
			// The context was cancelled, cancel the delay and return the error
			return action.FromError(ctx.Err())
		}
	}
}
