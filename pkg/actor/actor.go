package actor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/loopable"
	"time"
)

type Actor struct {
	actions   []action.Action
	loopDelay time.Duration
}

func (a *Actor) Do(ctx context.Context) (action.Result, error) {
	return loopable.RunForever(ctx, a.loopDelay, a.actions...)
}
