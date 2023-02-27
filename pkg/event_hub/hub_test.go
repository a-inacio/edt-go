package event_hub

import (
	"context"
	"errors"
	"fmt"
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

func (h *SomeEventHandler) Handler(ctx context.Context, event interface{}) error {
	h.GotCalled = true

	if event.(SomeEvent).ShouldFail {
		return errors.New("I was asked to fail")
	}

	return nil
}

type SomeGenericEventHandler struct {
	GotCalled bool
	Message   string
}

func (h *SomeGenericEventHandler) Handler(ctx context.Context, e interface{}) error {
	h.GotCalled = true

	values := e.(event.GenericNamedEvent).Values

	h.Message = fmt.Sprintf("%v", values["Message"])

	return nil
}

func TestHub_PublishAndSubscribe(t *testing.T) {
	hub := NewHub(nil)

	someEventHandler := &SomeEventHandler{}

	hub.Subscribe(SomeEvent{}, someEventHandler)

	wg := hub.Publish(SomeEvent{}, nil)

	wg.Wait()

	if !someEventHandler.GotCalled {
		t.Errorf("The handler should have been called")
	}
}

func TestHub_PublishEventWithoutSubscribers(t *testing.T) {
	hub := NewHub(nil)

	wg := hub.Publish(SomeOtherEvent{}, nil)

	wg.Wait()
}

func TestHub_FailingSubscriber(t *testing.T) {
	hub := NewHub(nil)

	someEventHandler := &SomeEventHandler{}

	hub.Subscribe(SomeEvent{}, someEventHandler)

	wg := hub.Publish(SomeEvent{ShouldFail: true}, nil)

	wg.Wait()

	if !someEventHandler.GotCalled {
		t.Errorf("The handler should have been called")
	}
}

func TestHub_PublishAndSubscribeWithGenericEvents(t *testing.T) {
	hub := NewHub(nil)

	someEventHandler := &SomeGenericEventHandler{}

	hub.Subscribe(event.WithName("SomeEvent"), someEventHandler)

	wg := hub.Publish(*event.WithNameAndKeyValues("SomeEvent", "Message", 42), nil)

	wg.Wait()

	if !someEventHandler.GotCalled {
		t.Errorf("The handler should have been called")
	}

	if someEventHandler.Message != "42" {
		t.Errorf("The handler should have received a message")
	}

}
