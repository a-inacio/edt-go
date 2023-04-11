package delayable

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Delayable struct {
	delay     time.Duration
	operation action.Action
}

type Builder struct {
	delay  time.Duration
	action action.Action
}
