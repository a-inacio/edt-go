package actor

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Actor struct {
	actions   []action.Action
	loopDelay time.Duration
}

type Builder struct {
	actions   []action.Action
	loopDelay time.Duration
}
