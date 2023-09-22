package expectable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/delayable"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/eventhub"
	"testing"
	"time"
)

type SomeEvent struct {
	Message string
}

func TestExpectable_ContinueAfterEvent(t *testing.T) {
	hub := eventhub.NewEventHub(nil)

	ctx := context.Background()

	expect := NewExpectable(hub, SomeEvent{})

	go delayable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{
			Message: "Hello EDT!",
		}, ctx)
		return action.Nothing()
	})

	res, err := expect.Do(ctx)

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
	hub := eventhub.NewEventHub(nil)

	ctx, cancel := context.WithCancel(context.Background())

	go delayable.RunAfter(nil, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		cancel()
		return action.Nothing()
	})

	_, err := NewExpectable(hub, SomeEvent{}).Do(ctx)

	if err == nil {
		t.Errorf("Should have been canceled and an error returned")
	}
}

func TestExpectableBuilder_ShouldBeCanceled(t *testing.T) {
	hub := eventhub.NewEventHub(nil)

	_, err := NewBuilder().
		On(hub).
		Expect(SomeEvent{}).
		WithTimeout(1 * time.Second).
		Do(context.Background())

	if err == nil {
		t.Errorf("Should have been canceled and an error returned")
	}
}

func TestExpectableBuilder_ShouldContinueAfterEvent(t *testing.T) {
	hub := eventhub.NewEventHub(nil)
	ctx := context.Background()

	expect, err := NewBuilder().
		On(hub).
		Expect(SomeEvent{}).
		WithTimeout(2 * time.Second).
		Build()

	go delayable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{
			Message: "Hello EDT!",
		}, ctx)
		return action.Nothing()
	})

	res, err := expect.Do(ctx)

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
	hub := eventhub.NewEventHub(nil)
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

	go delayable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
		hub.Publish(SomeEvent{
			Message: "Hello EDT!",
		}, ctx)
		return action.Nothing()
	})

	_, err = expect.Do(ctx)

	if err == nil {
		t.Errorf("Should have been canceled and an error returned")
	}
}
