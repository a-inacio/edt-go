package promisse

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/expirable"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	promise := Future(nil, expirable.NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			return 42, nil
		}).
		WithTimeout(2*time.Second).
		Do)

	res, err := ValueOf[int](promise)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}

	res, err = ValueOf[int](promise)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}

func TestFutureChain(t *testing.T) {
	promise := Future(nil,
		func(ctx context.Context) (action.Result, error) {
			return 20, nil
		}).
		Then(func(ctx context.Context) (action.Result, error) {
			chained, _ := FromContext[int](ctx)
			return *chained + 22, nil
		})

	res, err := ValueOf[int](promise)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}

	res, err = ValueOf[int](promise)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}
