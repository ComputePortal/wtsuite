package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
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
