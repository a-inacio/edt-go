package awaitable

import (
	"context"
	"time"
)

func RunAfter(ctx context.Context, timeout time.Duration, action func(ctx context.Context) (any, error)) (interface{}, error) {
	select {
	case <-time.After(timeout):
		// Wait for a certain duration
	case <-ctx.Done():
		// The context was cancelled, cancel the delay and return the error
		return nil, ctx.Err()
	}

	return action(ctx)
}
