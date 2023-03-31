package expectable

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
)

type Expectable struct {
	e event.Event
	h *event_hub.Hub
}

type expectableEventHandler struct {
	ch chan struct {
		action.Result
		error
	}
}
