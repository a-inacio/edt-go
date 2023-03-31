package injector

import (
	"testing"
)

type SomeValue struct {
	message string
	counter int
}

func TestWithContext_Singleton(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(SomeValue{message: "Hello EDT!"})

	value, err := GetValue(injector, SomeValue{})

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}
}

func TestWithContext_Factory(t *testing.T) {
	injector := WithContext(nil)

	counter := 0
	injector.SetFactory(func() SomeValue {
		counter++
		return SomeValue{counter: counter}
	})

	value, err := GetValue(injector, SomeValue{})

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.counter != 1 {
		t.Errorf("Expected %v, got %v", 1, value.counter)
	}

	value, _ = GetValue(injector, SomeValue{})

	if value.counter != 2 {
		t.Errorf("Expected %v, got %v", 2, value.counter)
	}
}
