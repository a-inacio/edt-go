package awaitable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

func RunAfter(ctx context.Context, timeout time.Duration, a action.Action) (action.Result, error) {
	select {
	case <-time.After(timeout):
		// Wait for a certain duration
	case <-ctx.Done():
		// The context was cancelled, cancel the delay and return the error
		return action.FromError(ctx.Err())
	}

	return a(ctx)
}
