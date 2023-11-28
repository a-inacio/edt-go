package event

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Event can be anything
type Event interface{}

// NamedEvent is anything that explicitly states it has a name for an event
type NamedEvent interface {
	EventName() string
}

// GetName returns the name of the event, for the special case of NameEvents
// the returned value is the result of the `EventName` function otherwise
// it falls back to the Type name (since anything can be an event).
// By design, the package name is not included.
func GetName(event Event) string {
	if n, ok := event.(NamedEvent); ok {
		return n.EventName()
	}

	if n, ok := event.(*NamedEvent); ok {
		return (*n).EventName()
	}

	t := reflect.TypeOf(event)

	// Here be dragons... check the test cases for a better understanding
	if reflect.PtrTo(t).Implements(reflect.TypeOf((*NamedEvent)(nil)).Elem()) {
		v := reflect.ValueOf(event)

		vp := reflect.New(t)
		vp.Elem().Set(v)

		if tt, ok := vp.Interface().(NamedEvent); ok {
			return tt.EventName()
		}
	}

	if t.Kind() == reflect.Ptr {
		// We don't want the `*` in case this is a pointer  otherwise if passing
		// by a pointer or value the events will be different.
		return strings.TrimPrefix(t.Elem().Name(), "*")
	}

	return t.Name()
}

// WithName creates a new GenericNamedEvent with the given name
func WithName(name string) *GenericNamedEvent {
	return &GenericNamedEvent{name: name}
}

// WithNameAndKeyValues creates a new GenericNamedEvent with the given name, initializing it with the given key values
func WithNameAndKeyValues(name string, kv ...interface{}) *GenericNamedEvent {
	m := make(map[string]interface{})
	for i := 0; i < len(kv); i += 2 {
		key := fmt.Sprintf("%v", kv[i])
		value := kv[i+1]
		m[key] = value
	}

	return &GenericNamedEvent{name: name, Values: m}
}

// ValueOf returns the value of the event as the given type, if the event cannot be converted to the given type an error is returned.
func ValueOf[T any](event Event) (*T, error) {
	// Cast the value to the desired type.
	typedVal, ok := event.(T)
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem()
		return nil, fmt.Errorf("value is not of type %s", t.String())
	}

	return &typedVal, nil
}

// FromContext attempts to retrieve an event of the given type from the context.
// If the event is not found, an error is returned.
// Use this method in callback functions associated to event subscriptions.
func FromContext[T any](ctx context.Context) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	i, ok := ctx.Value(t.PkgPath()).(T)

	if !ok {
		return nil, errors.New("could not get value from context")
	}
	return &i, nil
}
