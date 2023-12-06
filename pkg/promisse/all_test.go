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
				chained, _ := ChainedValueOf[int](ctx)
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
			if *res[0] != 42 {
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
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 1, nil // 11
			},
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 2, nil // 12
			},
		).
		Wait().
		Then(func(ctx context.Context) (action.Result, error) {
			chained, _ := ChainedSliceOf[int](ctx)     // [11, 12]
			return 19 + *chained[0] + *chained[1], nil // 42 = 19 + 11 + 12
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

func TestFutureAllWithError(t *testing.T) {
	gotCalled := false

	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 10, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 1, nil // 11
			},

			action.DoNothing,

			func(ctx context.Context) (action.Result, error) {
				return action.FromErrorf("I'm not up to it")
			},
		).
		Wait().
		Then(func(ctx context.Context) (action.Result, error) {
			// This should never be executed!
			gotCalled = true
			return action.Nothing()
		})

	_, err := promise.Do(nil)

	if err == nil {
		t.Errorf("Should have failed!")
	}

	if gotCalled {
		t.Errorf("Should not have been called!")
	}
}

func TestFutureAllWithBailout(t *testing.T) {
	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 10, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 1, nil // 11
			},
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 2, nil // 12
			},
		).
		WaitWithBailout().
		Then(func(ctx context.Context) (action.Result, error) {
			chained, _ := ChainedSliceOf[int](ctx)     // [11, 12]
			return 19 + *chained[0] + *chained[1], nil // 42 = 19 + 11 + 12
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

func TestFutureAllWithBailoutAndError(t *testing.T) {
	gotCalled := false

	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 10, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 1, nil // 11
			},

			action.DoNothing,

			func(ctx context.Context) (action.Result, error) {
				return action.FromErrorf("I'm not up to it")
			},
		).
		WaitWithBailout().
		Then(func(ctx context.Context) (action.Result, error) {
			// This should never be executed!
			gotCalled = true
			return action.Nothing()
		})

	_, err := promise.Do(nil)

	if err == nil {
		t.Errorf("Should have failed!")
	}

	if gotCalled {
		t.Errorf("Should not have been called!")
	}
}

func TestFutureAllWithCancel(t *testing.T) {
	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 10, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 1, nil // 11
			},
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 2, nil // 12
			},
		).
		WaitWithCancel().
		Then(func(ctx context.Context) (action.Result, error) {
			chained, _ := ChainedSliceOf[int](ctx)     // [11, 12]
			return 19 + *chained[0] + *chained[1], nil // 42 = 19 + 11 + 12
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

func TestFutureAllWithCancelAndError(t *testing.T) {
	gotCalled := false

	promise := Future(
		func(ctx context.Context) (action.Result, error) {
			return 10, nil
		}).
		All(
			func(ctx context.Context) (action.Result, error) {
				chained, _ := ChainedValueOf[int](ctx)
				return *chained + 1, nil // 11
			},

			action.DoNothing,

			func(ctx context.Context) (action.Result, error) {
				return action.FromErrorf("I'm not up to it")
			},
		).
		WaitWithCancel().
		Then(func(ctx context.Context) (action.Result, error) {
			// This should never be executed!
			gotCalled = true
			return action.Nothing()
		})

	_, err := promise.Do(nil)

	if err == nil {
		t.Errorf("Should have failed!")
	}

	if gotCalled {
		t.Errorf("Should not have been called!")
	}
}
