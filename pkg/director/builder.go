package director

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/director/breaker"
	"sync"
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) Launch(actions ...action.Action) *Builder {
	builder.actions = actions
	return builder
}

func (builder *Builder) BreakWith(breaker breaker.Breaker) *Builder {
	builder.breaker = breaker
	return builder
}

func (builder *Builder) Build() *Director {
	var wg sync.WaitGroup
	wg.Add(len(builder.actions))

	return &Director{
		actions: builder.actions,
		wg:      wg,
		breaker: builder.breaker,
	}
}

func (builder *Builder) Go(ctx context.Context) (action.Result, error) {
	return builder.Build().Go(ctx)
}
