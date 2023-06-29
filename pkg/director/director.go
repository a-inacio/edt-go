package director

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/director/breaker/contextbreaker"
)

func (d *Director) Go(ctx context.Context) (action.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	breaker := d.breaker
	if breaker == nil {
		breaker = contextbreaker.FromContext(ctx)
	}

	defer breaker.Release()

	for _, a := range d.actions {
		go func(ctx context.Context, a action.Action) {
			defer d.wg.Done()
			a(ctx)
		}(breaker.Context(), a)
	}

	breaker.Wait()

	// Wait for all actions to complete.
	d.wg.Wait()

	return action.Nothing()
}
