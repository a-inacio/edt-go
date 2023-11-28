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
	SomeValue  string
}

type SomeOtherEvent struct {
}

type SomeEventHandler struct {
	GotCalled bool
}

func (h *SomeEventHandler) Handler(ctx context.Context, e event.Event) error {
	h.GotCalled = true
	se, _ := event.ValueOf[SomeEvent](e)

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
	gne, _ := event.ValueOf[event.GenericNamedEvent](e)

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
	ev := event.WithName("SomeEvent")
	hub.RegisterHandler(ev, ToHandler(ev, func(ctx context.Context, e event.Event) error {
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

func TestHub_PublishAndSubscribeWithActionWithContext(t *testing.T) {
	hub := NewEventHub(nil)

	gotCalled := false
	result := ""

	hub.Subscribe(SomeEvent{}, func(ctx context.Context) (action.Result, error) {
		gotCalled = true
		ev, _ := event.FromContext[SomeEvent](ctx)
		result = ev.SomeValue
		return action.Nothing()
	})

	wg := hub.Publish(SomeEvent{SomeValue: "42"}, nil)

	wg.Wait()

	if !gotCalled {
		t.Errorf("The handler should have been called")
	}

	if result != "42" {
		t.Errorf("Should have received 42")
	}
}

func TestHub_UnregisterWithAction(t *testing.T) {
	hub := NewEventHub(nil)

	gotCalledCount := 0
	handler := hub.Subscribe(SomeEvent{}, func(ctx context.Context) (action.Result, error) {
		gotCalledCount += 1
		return action.Nothing()
	})

	wg := hub.Publish(SomeEvent{}, nil)
	wg.Wait()

	if gotCalledCount != 1 {
		t.Errorf("The callback should have been invoked the first time")
	}

	hub.UnregisterHandler(SomeEvent{}, handler)
	wg = hub.Publish(SomeEvent{}, nil)
	wg.Wait()

	if gotCalledCount > 1 {
		t.Errorf("The callback should not have been invoked this time")
	}
}

func TestHub_UnsubscribeWithAction(t *testing.T) {
	hub := NewEventHub(nil)

	gotCalledCount := 0
	handler := hub.Subscribe(SomeEvent{}, func(ctx context.Context) (action.Result, error) {
		gotCalledCount += 1
		return action.Nothing()
	})

	wg := hub.Publish(SomeEvent{}, nil)
	wg.Wait()

	if gotCalledCount != 1 {
		t.Errorf("The callback should have been invoked the first time")
	}

	hub.Unsubscribe(handler)
	wg = hub.Publish(SomeEvent{}, nil)
	wg.Wait()

	if gotCalledCount > 1 {
		t.Errorf("The callback should not have been invoked this time")
	}
}
