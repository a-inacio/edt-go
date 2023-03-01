package action

import "context"

type Result interface{}
type Action func(ctx context.Context) (Result, error)
