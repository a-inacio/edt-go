package loopable

import (
	"context"
	"time"
)

func RunForever(ctx context.Context, delay time.Duration, action func(ctx context.Context) (any, error)) (interface{}, error) {
	for {
		action(ctx)

		select {
		case <-time.After(delay):
			// Wait for a certain delay
		case <-ctx.Done():
			// The context was cancelled, cancel the delay and return the error
			return nil, ctx.Err()
		}
	}
}
