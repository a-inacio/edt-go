package executor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

// Executor is a simple implementation of action.Action that executes a list of actions in sequence.
// It can be utilised as a command queue for executing actions.
type Executor struct {
	mu      sync.Mutex
	actions []action.Action
}

// NewExecutor creates a new Executor.
func NewExecutor() *Executor {
	return &Executor{}
}

// Add adds an action to the list of actions to be executed.
func (e *Executor) Add(a action.Action) *Executor {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.actions = append(e.actions, a)

	return e
}

// Do executes one action of the list of actions. If there are no actions left, it behaves like action.DoNothing.
func (e *Executor) Do(ctx context.Context) (action.Result, error) {
	cb := action.DoNothing
	e.mu.Lock()

	if len(e.actions) > 0 {
		cb = e.actions[0]
		e.actions = e.actions[1:]
	}

	e.mu.Unlock()

	return cb(ctx)
}
