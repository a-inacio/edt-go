package loopable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"testing"
	"time"
)

func TestLoopable_Go(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	counterA := 0
	counterB := 0
	NewBuilder().
		WithDelay(200*time.Millisecond).
		LoopOn(
			func(ctx context.Context) (action.Result, error) {
				counterA++
				return action.Nothing()
			},
			action.DoNothing,
			func(ctx context.Context) (action.Result, error) {
				counterB++
				return action.Nothing()
			}).
		Build().
		Do(ctx)

	if counterA < 3 {
		t.Errorf("Expected >= 3, got %v", counterA)
	}

	if counterB < 3 {
		t.Errorf("Expected >= 3, got %v", counterB)
	}
}
