package promisse

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"testing"
)

func TestFutureAllSimple(t *testing.T) {
	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 20, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := FromContext[int](ctx)
				return *chained + 22, nil
			},
		).
		Wait()

	go promise.Do(nil)

	res, err := SliceOf[int](promise)

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	} else {
		if len(res) != 1 {
			t.Errorf("Expected 1 element, got %v", len(res))
		} else {
			if res[0] != 42 {
				t.Errorf("Expected 42, got %v", res[0])
			}
		}
	}
}

func TestFutureAll(t *testing.T) {
	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 10, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := FromContext[int](ctx)
				return *chained + 1, nil // 11
			},
			func(ctx context.Context) (action.Result, error) {
				chained, _ := FromContext[int](ctx)
				return *chained + 2, nil // 12
			},
		).
		Wait().
		Then(func(ctx context.Context) (action.Result, error) {
			chained, _ := SliceFromContext[int](ctx) // [11, 12]
			return 19 + chained[0] + chained[1], nil // 42 = 19 + 11 + 12
		})

	go promise.Do(nil)

	res, err := ValueOf[int](promise)

	if *res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}
}
