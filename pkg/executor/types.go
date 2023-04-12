package executor

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

type Executor struct {
	mu      sync.Mutex
	actions []action.Action
}
