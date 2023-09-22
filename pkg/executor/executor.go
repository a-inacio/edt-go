package executor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

type Executor struct {
	mu      sync.Mutex
	actions []action.Action
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Add(a action.Action) *Executor {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.actions = append(e.actions, a)

	return e
}

func (e *Executor) ExecuteOne(ctx context.Context) (action.Result, error) {
	cb := action.DoNothing
	e.mu.Lock()

	if len(e.actions) > 0 {
		cb = e.actions[0]
		e.actions = e.actions[1:]
	}

	e.mu.Unlock()

	return cb(ctx)
}
