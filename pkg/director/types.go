package director

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/director/breaker"
	"sync"
)

type Director struct {
	actions []action.Action
	wg      sync.WaitGroup
	breaker breaker.Breaker
}

type Builder struct {
	actions []action.Action
	breaker breaker.Breaker
}
