package expectable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"github.com/a-inacio/edt-go/pkg/event_hub"
	"testing"
	"time"
)

type SomeEvent struct {
	ShouldFail bool
}

func TestExpectable_ContinueAfterEvent(t *testing.T) {
	hub := event_hub.NewHub(nil)

	ctx := context.Background()

	expect := NewExpectable(hub, SomeEvent{})

	go awaitable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{}, ctx)
		return action.Nothing()
	})

	expect.Go(ctx)
}

func TestExpectable_ShouldBeCanceled(t *testing.T) {
	hub := event_hub.NewHub(nil)

	ctx, cancel := context.WithCancel(context.Background())

	go awaitable.RunAfter(nil, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		cancel()
		return action.Nothing()
	})

	_, err := NewExpectable(hub, SomeEvent{}).Go(ctx)

	if err == nil {
		t.Errorf("Should have been canceled and an error returned")
	}
}
