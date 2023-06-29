package contextbreaker

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/director/breaker"
)

type impl struct {
	ctx context.Context
}

func FromContext(ctx context.Context) breaker.Breaker {
	return &impl{
		ctx: ctx,
	}
}

func (b *impl) Context() context.Context {
	return b.ctx
}

func (b *impl) Release() {
}

func (b *impl) Wait() {
	<-b.ctx.Done()
}
