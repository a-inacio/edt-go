package executor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"testing"
)

func TestExecutor_ExecuteOne(t *testing.T) {
	ex := NewExecutor().
		Add(func(ctx context.Context) (action.Result, error) {
			return "A", nil
		}).
		Add(func(ctx context.Context) (action.Result, error) {
			return 1, nil
		}).
		Add(func(ctx context.Context) (action.Result, error) {
			return true, nil
		})

	v1, _ := ex.Do(nil)
	v2, _ := ex.Do(nil)
	v3, _ := ex.Do(nil)
	v4, _ := ex.Do(nil)

	s1, _ := action.ValueOf[string](v1)
	i2, _ := action.ValueOf[int](v2)
	b3, _ := action.ValueOf[bool](v3)

	if *s1 != "A" {
		t.Errorf("expected A, got %v", s1)
	}

	if *i2 != 1 {
		t.Errorf("expected 1, got %v", i2)
	}

	if *b3 != true {
		t.Errorf("expected true, got %v", b3)
	}

	if v4 != nil {
		t.Errorf("expected nil, got %v", v4)
	}
}
