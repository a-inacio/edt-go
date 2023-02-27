package event

import "reflect"

// Event can be anything
type Event interface{}

// NamedEvent is anything that explicitly states it has a name for an event
type NamedEvent interface {
	EventName() string
}

func GetName(event Event) string {
	return reflect.TypeOf(event).Name()
}
