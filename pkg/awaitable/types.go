package awaitable

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

type Awaitable struct {
	ctx context.Context
	wg  sync.WaitGroup
	r   action.Result
	e   error
}
