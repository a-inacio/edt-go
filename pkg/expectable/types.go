package expectable

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/eventhub"
	"time"
)

type Expectable struct {
	e        event.Event
	h        *eventhub.EventHub
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
	h        *eventhub.EventHub
	timeout  time.Duration
	criteria func(e event.Event) bool
}
