package actor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/loopable"
)

func (a *Actor) Go(ctx context.Context) (action.Result, error) {
	return loopable.RunForever(ctx, a.loopDelay, a.actions...)
}
