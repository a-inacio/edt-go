package statemachine

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/event"
)

type stateMachineHubHandler struct {
	sm *StateMachine
}

func (h *stateMachineHubHandler) Handler(ctx context.Context, e event.Event) error {
	return h.sm.TriggerEvent(e)
}
