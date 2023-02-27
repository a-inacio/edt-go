package state_machine

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
)

func (h *stateMachineHubHandler) Handler(ctx context.Context, e event.Event) error {
	return h.sm.TriggerEvent(e)
}
