package expectable

import (
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
	"sync"
)

type Expectable struct {
	e event.Event
	h *event_hub.Hub
}

type expectableEventHandler struct {
	wg sync.WaitGroup
}
