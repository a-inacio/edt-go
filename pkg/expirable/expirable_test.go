package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/delayable"
	"testing"
	"time"
)

func TestExpirable_Expect42(t *testing.T) {
	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			return 42, nil
		}).
		WithTimeout(2 * time.Second).
		Do(nil)

	if res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}

func TestExpirable_ExpectTimeout(t *testing.T) {
	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			return delayable.RunAfter(ctx, 5*time.Second, func(ctx context.Context) (action.Result, error) {
				return 42, nil
			})
		}).
		WithTimeout(1 * time.Second).
		Do(nil)

	if res == 42 {
		t.Errorf("Should not have got 42")
	}

	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestExpirable_ShouldNotTimeout(t *testing.T) {
	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			return delayable.RunAfter(ctx, 1*time.Second, func(ctx context.Context) (action.Result, error) {
				return 42, nil
			})
		}).
		WithTimeout(2 * time.Second).
		Do(nil)

	if res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}

func TestExpirable_ExpectCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	operationCalled := false

	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			operationCalled = true

			return delayable.RunAfter(ctx, 5*time.Second, func(ctx context.Context) (action.Result, error) {
				return 42, nil
			})
		}).
		WithTimeout(10 * time.Second).
		Do(ctx)

	if !operationCalled {
		t.Errorf("action should have been called")
	}

	if res == 42 {
		t.Errorf("Should not have got 42")
	}

	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestExpirable_ExpectCancellationDuringDelay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	operationCalled := false

	res, err := NewBuilder().
		WithDelay(2 * time.Second).
		FromOperation(func(ctx context.Context) (action.Result, error) {
			operationCalled = true

			return delayable.RunAfter(ctx, 5*time.Second, func(ctx context.Context) (action.Result, error) {
				return 42, nil
			})
		}).
		WithTimeout(10 * time.Second).
		Do(ctx)

	if operationCalled {
		t.Errorf("Action should not have been called")
	}

	if res == 42 {
		t.Errorf("Should not have got 42")
	}

	if err == nil {
		t.Errorf("Should have failed")
	}
}
