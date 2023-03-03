package director

import (
	"github.com/a-inacio/edt-go/pkg/action"
	"sync"
)

type Director struct {
	actions []action.Action
	wg      sync.WaitGroup
}

type Builder struct {
	actions []action.Action
}
