package event

import (
	"fmt"
	"reflect"
	"strings"
)

// GetName returns the name of the even, for the special case of NameEvents
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

func WithName(name string) *GenericNamedEvent {
	return &GenericNamedEvent{name: name}
}

func WithNameAndKeyValues(name string, kv ...interface{}) *GenericNamedEvent {
	m := make(map[string]interface{})
	for i := 0; i < len(kv); i += 2 {
		key := fmt.Sprintf("%v", kv[i])
		value := kv[i+1]
		m[key] = value
	}

	return &GenericNamedEvent{name: name, Values: m}
}
