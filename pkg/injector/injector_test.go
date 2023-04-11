package injector

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/expirable"
	"testing"
	"time"
)

type SomeValue struct {
	message string
	counter int
}

type AnotherValue struct {
	message string
	counter int
}

func TestWithContext_Singleton(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(SomeValue{message: "Hello EDT!"})

	value, err := GetValue[SomeValue](injector)

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

	value, err := GetValue[SomeValue](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.counter != 1 {
		t.Errorf("Expected %v, got %v", 1, value.counter)
	}

	value, _ = GetValue[SomeValue](injector)

	if value.counter != 2 {
		t.Errorf("Expected %v, got %v", 2, value.counter)
	}
}

func TestFromContext(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(SomeValue{message: "Hello EDT!"})

	value, err := expirable.NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			value, err := GetValue[SomeValue](FromContext(ctx))
			return value.message, err
		}).
		WithTimeout(2 * time.Second).
		Go(injector.Context())

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value)
	}
}

func TestChainMethods(t *testing.T) {
	counter := 0

	ctx := WithContext(nil).
		SetSingleton(SomeValue{message: "Hello EDT!"}).
		SetFactory(func() AnotherValue {
			counter++
			return AnotherValue{counter: counter}
		}).
		Context()

	value, err := GetValueFromContext[SomeValue](ctx)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}

	anotherValue, err := GetValueFromContext[AnotherValue](ctx)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if anotherValue.counter != 1 {
		t.Errorf("Expected %v, got %v", 1, value.counter)
	}
}
