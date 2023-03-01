package expirable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"time"
)

type Hooks struct {
	Action      func(ctx context.Context) (interface{}, error)
	OnExpired   func(ctx context.Context)
	OnSuccess   func(ctx context.Context, result interface{})
	OnError     func(ctx context.Context, e error)
	OnCanceled  func(ctx context.Context)
	OnCompleted func(ctx context.Context, result interface{}, err error)
}

type Expirable struct {
	timeout time.Duration
	action  action.Action
	hooks   Hooks
}

type Builder struct {
	timeout   time.Duration
	delay     time.Duration
	operation action.Action
	hooks     Hooks
}
