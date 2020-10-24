package js

import (
	"../context"
)

type Module interface {
	GetExportedVariable(gs GlobalScope, name string,
		nameCtx context.Context) (Variable, error)

	Context() context.Context
}
