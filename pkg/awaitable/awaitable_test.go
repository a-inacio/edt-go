package awaitable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/expirable"
	"testing"
	"time"
)

func TestAwaitFor(t *testing.T) {
	awaitable := AwaitFor(nil, expirable.NewBuilder().
		FromOperation(func(ctx context.Context) (action.Result, error) {
			return 42, nil
		}).
		WithTimeout(2*time.Second).
		Do)

	res, err := GetValue[int](awaitable)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}

	res, err = GetValue[int](awaitable)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}
