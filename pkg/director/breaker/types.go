package breaker

import (
	"context"
)

type Breaker interface {
	Context() context.Context
	Release()
	Wait()
}
