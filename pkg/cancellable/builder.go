package cancellable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

type Builder struct {
	action action.Action
	wg     *sync.WaitGroup
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) FromAction(action action.Action) *Builder {
	builder.action = action
	return builder
}

func (builder *Builder) WithWaitGroup(wg sync.WaitGroup) *Builder {
	builder.wg = &wg
	return builder
}

func (builder *Builder) Build() *Cancellable {
	wg := builder.wg
	if wg == nil {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	}
	inst := &Cancellable{
		action: builder.action,
		wg:     wg,
	}

	inst.started.Add(1)

	return inst
}

func (builder *Builder) Do(ctx context.Context) (action.Result, error) {
	return builder.Build().Do(ctx)
}
