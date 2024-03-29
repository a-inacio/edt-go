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

type SomeInterface interface {
	SomeMethod() string
}

type SomeTypeWithInterface struct {
	message string
}

func (t SomeTypeWithInterface) SomeMethod() string { return t.message }

type YetAnotherTypeWithInterface struct {
	message string
}

func NewYetAnotherTypeWithInterface(message string) SomeInterface {
	return YetAnotherTypeWithInterface{
		message: message,
	}
}

func NewYetAnotherTypeWithInterfacePtr(message string) SomeInterface {
	return &YetAnotherTypeWithInterface{
		message: message,
	}
}

func (t YetAnotherTypeWithInterface) SomeMethod() string { return t.message }

func TestWithContext_Singleton(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(SomeValue{message: "Hello EDT!"})

	value, err := Get[SomeValue](injector)

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

func TestWithContext_Singleton_Func(t *testing.T) {
	injector := WithContext(nil)

	counter := 0
	injector.SetSingleton(func() SomeValue {
		counter++
		return SomeValue{message: "Hello EDT!"}
	})

	value, err := Get[SomeValue](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}

	Get[SomeValue](injector)

	if counter != 1 {
		t.Errorf("Should have been called once, expected %v, got %v", 1, counter)
	}
}

func TestWithContext_Singleton_Func_WithDependencies(t *testing.T) {
	injector := WithContext(nil)

	counter := 0
	injector.SetSingleton(func() SomeValue {
		counter++
		return SomeValue{message: "Hello EDT!"}
	})

	anotherCounter := 0
	injector.SetSingleton(func(value SomeValue) AnotherValue {
		anotherCounter++
		return AnotherValue{message: value.message}
	})

	value, err := Get[AnotherValue](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}

	Get[AnotherValue](injector)

	if counter != 1 {
		t.Errorf("SomeValue singleton should have been called once, expected %v, got %v", 1, counter)
	}

	if anotherCounter != 1 {
		t.Errorf("AnotherValue singleton should have been called once, expected %v, got %v", 1, anotherCounter)
	}
}

func TestWithContext_Satisfy_Func(t *testing.T) {
	injector := WithContext(nil)

	injector.SetSingleton(func() SomeValue {
		return SomeValue{message: "Hello EDT!"}
	})

	value, err := Resolve[AnotherValue](injector, func(value SomeValue) AnotherValue {
		return AnotherValue{message: value.message}
	})

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

func TestWithContext_MustSatisfy_Func(t *testing.T) {
	injector := WithContext(nil)

	injector.SetSingleton(func() SomeValue {
		return SomeValue{message: "Hello EDT!"}
	})

	value := MustResolve[AnotherValue](injector, func(value SomeValue) AnotherValue {
		return AnotherValue{message: value.message}
	})

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

	value, err := Get[SomeValue](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.counter != 1 {
		t.Errorf("Expected %v, got %v", 1, value.counter)
	}

	value, _ = Get[SomeValue](injector)

	if value.counter != 2 {
		t.Errorf("Expected %v, got %v", 2, value.counter)
	}
}

func TestWithContext_FactoryWithArguments(t *testing.T) {
	injector := WithContext(nil)

	counter := 0
	injector.SetSingleton(AnotherValue{counter: 10})
	injector.SetFactory(func(value AnotherValue) SomeValue {
		counter++
		return SomeValue{counter: value.counter + counter}
	})

	value, err := Get[SomeValue](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if value.counter != 11 {
		t.Errorf("Expected %v, got %v", 11, value.counter)
	}

	value, _ = Get[SomeValue](injector)

	if value.counter != 12 {
		t.Errorf("Expected %v, got %v", 12, value.counter)
	}
}

func TestFromContext(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(SomeValue{message: "Hello EDT!"})

	value, err := expirable.NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			value, err := Get[SomeValue](FromContext(ctx))
			return value.message, err
		}).
		WithTimeout(2 * time.Second).
		Do(injector.Context())

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

	value, err := GetFromContext[SomeValue](ctx)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}

	anotherValue, err := GetFromContext[AnotherValue](ctx)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if anotherValue.counter != 1 {
		t.Errorf("Expected %v, got %v", 1, value.counter)
	}
}

func TestWithContext_Singleton_WithInterface(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(SomeTypeWithInterface{message: "Hello EDT!"})

	value, err := Get[SomeInterface](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if (*value).SomeMethod() != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", (*value).SomeMethod())
	}
}

func TestWithContext_Singleton_Ptr(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(&SomeValue{message: "Hello EDT!"})

	value, err := Get[SomeValue](injector)

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

func TestWithContext_Singleton_WithNoSingleInterface(t *testing.T) {
	injector := WithContext(nil)
	injector.
		SetSingleton(SomeTypeWithInterface{message: "Hello EDT!"}).
		SetSingleton(YetAnotherTypeWithInterface{message: "Hello EDT!"})

	_, err := Get[SomeInterface](injector)

	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestWithContext_Singleton_WithInterfaceAndConstructorWithInterface(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(NewYetAnotherTypeWithInterface("Hello EDT!"))

	value, err := Get[SomeInterface](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if (*value).SomeMethod() != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", (*value).SomeMethod())
	}
}

func TestWithContext_Singleton_WithInterfaceAndConstructorWithInterfacePtr(t *testing.T) {
	injector := WithContext(nil)
	injector.SetSingleton(NewYetAnotherTypeWithInterfacePtr("Hello EDT!"))

	value, err := Get[SomeInterface](injector)

	if err != nil {
		t.Errorf("Should not have failed")
	}

	if value == nil {
		t.Errorf("Should have gotten a value")
	}

	if (*value).SomeMethod() != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", (*value).SomeMethod())
	}
}

func TestWithContext_MustSatisfy_Interface(t *testing.T) {
	injector := WithContext(nil)

	injector.
		SetSingleton(NewYetAnotherTypeWithInterface("Hello EDT!")).
		SetSingleton(&SomeValue{message: "Hello EDT!"})

	value := MustResolve[AnotherValue](injector, func(value SomeInterface) AnotherValue {
		return AnotherValue{message: value.SomeMethod()}
	})

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}

	anotherValue := MustResolve[SomeInterface](injector, func(value *SomeValue) SomeInterface {
		return NewYetAnotherTypeWithInterface(value.message)
	})

	if anotherValue == nil {
		t.Errorf("Should have gotten a value")
	}

	if anotherValue.SomeMethod() != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", anotherValue.SomeMethod())
	}
}

func TestWithContext_MustSatisfy_InterfacePtr(t *testing.T) {
	injector := WithContext(nil)

	injector.
		SetSingleton(NewYetAnotherTypeWithInterfacePtr("Hello EDT!")).
		SetSingleton(&SomeValue{message: "Hello EDT!"})

	value := MustResolve[AnotherValue](injector, func(value SomeInterface) AnotherValue {
		return AnotherValue{message: value.SomeMethod()}
	})

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}

	anotherValue := MustResolve[SomeInterface](injector, func(value *SomeValue) SomeInterface {
		return NewYetAnotherTypeWithInterfacePtr(value.message)
	})

	if anotherValue.SomeMethod() != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", anotherValue.SomeMethod())
	}
}

func TestMustGetValue(t *testing.T) {
	injector := WithContext(nil)

	injector.
		SetSingleton(NewYetAnotherTypeWithInterfacePtr("Hello EDT!")).
		SetSingleton(&SomeValue{message: "Hello EDT!"})

	value := MustGet[SomeValue](injector)

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}
}

func TestMustGetValueFromContext(t *testing.T) {
	injector := WithContext(nil)

	injector.
		SetSingleton(NewYetAnotherTypeWithInterfacePtr("Hello EDT!")).
		SetSingleton(&SomeValue{message: "Hello EDT!"})

	value := MustGetFromContext[SomeValue](injector.Context())

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}
}

func TestMustResolveFromContext_Func(t *testing.T) {
	ctx := WithContext(nil).
		SetSingleton(func() SomeValue {
			return SomeValue{message: "Hello EDT!"}
		}).
		Context()

	value := MustResolveFromContext[AnotherValue](ctx, func(value SomeValue) AnotherValue {
		return AnotherValue{message: value.message}
	})

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}
}

func TestSatisfyContextWhenRequired(t *testing.T) {
	injector := WithContext(nil)

	ctx := injector.Context()

	value := MustGet[context.Context](injector)

	if ctx != value {
		t.Errorf("Expected to be able to resolve a context.Context")
	}

	value = MustResolve[context.Context](injector, func(ctx context.Context) context.Context {
		return ctx
	})

	if ctx != value {
		t.Errorf("Expected to be able to resolve a context.Context")
	}
}

func TestMustGetValuePtr(t *testing.T) {
	injector := WithContext(nil)

	injector.
		SetSingleton(NewYetAnotherTypeWithInterfacePtr("Hello EDT!")).
		SetSingleton(&SomeValue{message: "Hello EDT!"})

	value := MustGet[*SomeValue](injector)

	if value.message != "Hello EDT!" {
		t.Errorf("Expected %s, got %s", "Hello EDT!", value.message)
	}
}
