package expirable

import (
	"context"
	"time"
)

type Hooks struct {
	Operation   func(ctx context.Context) (interface{}, error)
	OnExpired   func(ctx context.Context)
	OnSuccess   func(ctx context.Context, result interface{})
	OnError     func(ctx context.Context, e error)
	OnCanceled  func(ctx context.Context)
	OnCompleted func(ctx context.Context, result interface{}, err error)
}

type Expirable struct {
	timeout   time.Duration
	operation func(ctx context.Context) (interface{}, error)
	hooks     Hooks
}

type Builder struct {
	timeout   time.Duration
	operation func(ctx context.Context) (interface{}, error)
	hooks     Hooks
}
