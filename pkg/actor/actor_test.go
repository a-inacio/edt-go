package actor

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"testing"
	"time"
)

func TestActor_Go(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	count := 0
	NewBuilder().
		LoopingForever(200*time.Millisecond,
			func(ctx context.Context) (action.Result, error) {
				count++
				return action.Nothing()
			},
			action.DoNothing,
			action.DoNothing).
		Do(ctx)

	if count < 4 || count > 5 {
		t.Errorf("Action executed an unpected number of times, %v", count)
	}
}
