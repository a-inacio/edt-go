package signalbreaker

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/director/breaker"
	"os/signal"
	"syscall"
)

type impl struct {
	ctx  context.Context
	done context.CancelFunc
}

func FromContext(ctx context.Context) breaker.Breaker {
	// Create a cancellation context
	dCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	return &impl{
		ctx:  dCtx,
		done: stop,
	}
}

func (b *impl) Context() context.Context {
	return b.ctx
}

func (b *impl) Release() {
	b.done()
}

func (b *impl) Wait() {
	<-b.ctx.Done()
	b.done()
}
