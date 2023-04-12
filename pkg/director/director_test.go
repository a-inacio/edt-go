package director

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/actor"
	"testing"
	"time"
)

func TestDirector_Go(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	counterA := 0
	counterB := 0
	NewBuilder().
		Launch(
			func(ctx context.Context) (action.Result, error) {
				counterA++
				return action.Nothing()
			},
			action.DoNothing,
			actor.NewBuilder().
				LoopingForever(200*time.Millisecond, func(ctx context.Context) (action.Result, error) {
					counterB++
					return action.Nothing()
				}).
				Build().
				Go).
		Build().
		Go(ctx)

	if counterA != 1 {
		t.Errorf("Expected 1, got %v", counterA)
	}

	if counterB < 3 {
		t.Errorf("Expected >= 3, got %v", counterB)
	}
}
