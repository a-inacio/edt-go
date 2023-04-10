package expectable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/event_hub"
	"testing"
	"time"
)

type SomeEvent struct {
	Message string
}

func TestExpectable_ContinueAfterEvent(t *testing.T) {
	hub := event_hub.NewHub(nil)

	ctx := context.Background()

	expect := NewExpectable(hub, SomeEvent{})

	go awaitable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{
			Message: "Hello EDT!",
		}, ctx)
		return action.Nothing()
	})

	res, err := expect.Go(ctx)

	if err != nil {
		t.Errorf("Should not have failled")
	}

	if res == nil {
		t.Errorf("Should have gotten a result")
	}

	if res.(SomeEvent).Message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", res.(SomeEvent).Message)
	}
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

func TestExpectableBuilder_ShouldBeCanceled(t *testing.T) {
	hub := event_hub.NewHub(nil)

	_, err := NewBuilder().
		On(hub).
		Expect(SomeEvent{}).
		WithTimeout(1 * time.Second).
		Go(context.Background())

	if err == nil {
		t.Errorf("Should have been canceled and an error returned")
	}
}

func TestExpectableBuilder_ShouldContinueAfterEvent(t *testing.T) {
	hub := event_hub.NewHub(nil)
	ctx := context.Background()

	expect, err := NewBuilder().
		On(hub).
		Expect(SomeEvent{}).
		WithTimeout(2 * time.Second).
		Build()

	go awaitable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{
			Message: "Hello EDT!",
		}, ctx)
		return action.Nothing()
	})

	res, err := expect.Go(ctx)

	if err != nil {
		t.Errorf("Should not have failled")
	}

	if res == nil {
		t.Errorf("Should have gotten a result")
	}

	if res.(SomeEvent).Message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", res.(SomeEvent).Message)
	}
}

func TestExpectableBuilder_ShouldNotContinueAfterEvent(t *testing.T) {
	hub := event_hub.NewHub(nil)
	ctx := context.Background()

	expect, err := NewBuilder().
		On(hub).
		Expect(SomeEvent{}).
		WithTimeout(2 * time.Second).
		Where(func(e event.Event) bool {
			if e.(SomeEvent).Message == "Hello EDT!" {
				return false
			}

			return true
		}).
		Build()

	go awaitable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{
			Message: "Hello EDT!",
		}, ctx)
		return action.Nothing()
	})

	_, err = expect.Go(ctx)

	if err == nil {
		t.Errorf("Should have been canceled and an error returned")
	}
}
