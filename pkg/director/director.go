package director

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"os"
	"os/signal"
	"syscall"
)

func (d *Director) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a cancellation context
	dCtx, cancel := context.WithCancel(context.Background())

	for _, a := range d.actions {
		go func(ctx context.Context, a action.Action) {
			defer d.wg.Done()
			a(ctx)
		}(dCtx, a)
	}

	// Wait for a signal to shut down the server
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	cancel()

	// Wait for all actions to complete.
	d.wg.Wait()

	return action.Nothing()
}
