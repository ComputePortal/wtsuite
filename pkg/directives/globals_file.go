package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

const FILE = "__file__"

func SetFile(scope Scope, path string, ctx context.Context) {
	// set the __file__ internal variable immediately
	scope.SetVar(FILE, functions.Var{
		tokens.NewValueString(path, ctx),
		true,
		true,
		false,
		false,
		ctx,
	})
}
