package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Module interface {
	GetExportedVariable(gs GlobalScope, name string,
		nameCtx context.Context) (Variable, error)

	Context() context.Context
}
