package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"testing"
	"time"
)

func TestExpirable_Expect42(t *testing.T) {
	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (interface{}, error) {
			return 42, nil
		}).
		WithTimeout(2 * time.Second).
		Go(nil)

	if res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}

func TestExpirable_ExpectTimeout(t *testing.T) {
	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (interface{}, error) {
			return awaitable.RunAfter(ctx, 5*time.Second, func() (any, error) {
				return 42, nil
			})
		}).
		WithTimeout(1 * time.Second).
		Go(nil)

	if res == 42 {
		t.Errorf("Should not have got 42")
	}

	if err == nil {
		t.Errorf("Should have failed")
	}
}

func TestExpirable_ShouldNotTimeout(t *testing.T) {
	res, err := NewBuilder().
		FromOperation(func(ctx context.Context) (interface{}, error) {
			return awaitable.RunAfter(ctx, 1*time.Second, func() (any, error) {
				return 42, nil
			})
		}).
		WithTimeout(2 * time.Second).
		Go(nil)

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
		FromOperation(func(ctx context.Context) (interface{}, error) {
			operationCalled = true

			return awaitable.RunAfter(ctx, 5*time.Second, func() (any, error) {
				return 42, nil
			})
		}).
		WithTimeout(10 * time.Second).
		Go(ctx)

	if !operationCalled {
		t.Errorf("operation should have been called")
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
		FromOperation(func(ctx context.Context) (interface{}, error) {
			operationCalled = true

			return awaitable.RunAfter(ctx, 5*time.Second, func() (any, error) {
				return 42, nil
			})
		}).
		WithTimeout(10 * time.Second).
		Go(ctx)

	if operationCalled {
		t.Errorf("Operation should not have been called")
	}

	if res == 42 {
		t.Errorf("Should not have got 42")
	}

	if err == nil {
		t.Errorf("Should have failed")
	}
}
