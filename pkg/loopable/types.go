package loopable

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Loopable struct {
	actions []action.Action
	delay   time.Duration
}

type Builder struct {
	actions []action.Action
	delay   time.Duration
}
