package cancellable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"testing"
)

func TestCancellable_Expect42(t *testing.T) {
	cancellable := NewBuilder().
		FromAction(func(ctx context.Context) (action.Result, error) {
			return 42, nil
		}).
		Build()

	res, err := cancellable.Do(nil)

	if res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}

	// Even if already finished, it should still be possible to wait for it.
	res, err = cancellable.Wait(nil)

	if res != 42 {
		t.Errorf("Expected 42, got %v", res)
	}

	if err != nil {
		t.Errorf("Should have not failed - %v", err)
	}

	// Even if already finished, it should still be possible to safely cancel it.
	cancellable.Cancel()
}
