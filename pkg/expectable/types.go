package expectable

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
	"time"
)

type Expectable struct {
	e        event.Event
	h        *event_hub.Hub
	timeout  time.Duration
	criteria func(e event.Event) bool
}

type expectableEventHandler struct {
	ch chan struct {
		action.Result
		error
	}
}

type Builder struct {
	e        event.Event
	h        *event_hub.Hub
	timeout  time.Duration
	criteria func(e event.Event) bool
}
