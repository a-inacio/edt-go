package eventhub

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/event"
	"testing"
)

type SomeEvent struct {
	ShouldFail bool
}

type SomeOtherEvent struct {
}

type SomeEventHandler struct {
	GotCalled bool
}

func (h *SomeEventHandler) Handler(ctx context.Context, e event.Event) error {
	h.GotCalled = true
	se, _ := event.GetValue[SomeEvent](e)

	if se.ShouldFail {
		return errors.New("I was asked to fail")
	}

	return nil
}

type SomeGenericEventHandler struct {
	GotCalled bool
	Message   string
}

func (h *SomeGenericEventHandler) Handler(ctx context.Context, e event.Event) error {
	h.GotCalled = true
	gne, _ := event.GetValue[event.GenericNamedEvent](e)

	values := gne.Values

	h.Message = fmt.Sprintf("%v", values["Message"])

	return nil
}

func TestHub_PublishAndSubscribe(t *testing.T) {
	hub := NewEventHub(nil)

	someEventHandler := &SomeEventHandler{}

	hub.RegisterHandler(SomeEvent{}, someEventHandler)

	wg := hub.Publish(SomeEvent{}, nil)

	wg.Wait()

	if !someEventHandler.GotCalled {
		t.Errorf("The handler should have been called")
	}
}

func TestHub_PublishEventWithoutSubscribers(t *testing.T) {
	hub := NewEventHub(nil)

	wg := hub.Publish(SomeOtherEvent{}, nil)

	wg.Wait()
}

func TestHub_FailingSubscriber(t *testing.T) {
	hub := NewEventHub(nil)

	someEventHandler := &SomeEventHandler{}

	hub.RegisterHandler(SomeEvent{}, someEventHandler)

	wg := hub.Publish(SomeEvent{ShouldFail: true}, nil)

	wg.Wait()

	if !someEventHandler.GotCalled {
		t.Errorf("The handler should have been called")
	}
}

func TestHub_PublishAndSubscribeWithGenericEvents(t *testing.T) {
	hub := NewEventHub(nil)

	someEventHandler := &SomeGenericEventHandler{}

	hub.RegisterHandler(event.WithName("SomeEvent"), someEventHandler)

	wg := hub.Publish(*event.WithNameAndKeyValues("SomeEvent", "Message", 42), nil)

	wg.Wait()

	if !someEventHandler.GotCalled {
		t.Errorf("The handler should have been called")
	}

	if someEventHandler.Message != "42" {
		t.Errorf("The handler should have received a message")
	}
}

func TestHub_PublishAndSubscribeWithGenericCallbacks(t *testing.T) {
	hub := NewEventHub(nil)

	gotCalled := 0
	hub.RegisterHandler(event.WithName("SomeEvent"), ToHandler(func(ctx context.Context, e event.Event) error {
		gotCalled++
		return nil
	}))

	wg := hub.Publish(*event.WithName("SomeEvent"), nil)

	wg.Wait()

	if gotCalled != 1 {
		t.Errorf("The handler should have been called once")
	}
}

func TestHub_PublishAndSubscribeWithAction(t *testing.T) {
	hub := NewEventHub(nil)

	gotCalled := false
	hub.Subscribe(SomeEvent{}, func(ctx context.Context) (action.Result, error) {
		gotCalled = true
		return action.Nothing()
	})

	wg := hub.Publish(SomeEvent{}, nil)

	wg.Wait()

	if !gotCalled {
		t.Errorf("The handler should have been called")
	}
}
