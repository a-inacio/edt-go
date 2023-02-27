package event

// Event can be anything
type Event interface{}

// NamedEvent is anything that explicitly states it has a name for an event
type NamedEvent interface {
	EventName() string
}

type GenericNamedEvent struct {
	name   string
	Values map[string]interface{}
}
