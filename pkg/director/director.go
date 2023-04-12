package director

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"os/signal"
	"syscall"
)

func (d *Director) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a cancellation context
	dCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for _, a := range d.actions {
		go func(ctx context.Context, a action.Action) {
			defer d.wg.Done()
			a(ctx)
		}(dCtx, a)
	}

	// Listen for the interrupt signal.
	<-dCtx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()

	// Wait for all actions to complete.
	d.wg.Wait()

	return action.Nothing()
}
