package event

type GenericNamedEvent struct {
	name   string
	Values map[string]interface{}
}

func (e *GenericNamedEvent) EventName() string {
	return e.name
}
