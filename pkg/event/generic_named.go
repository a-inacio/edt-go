package event

func (e *GenericNamedEvent) EventName() string {
	return e.name
}
