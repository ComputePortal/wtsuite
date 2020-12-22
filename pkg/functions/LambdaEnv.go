package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type LambdaScope interface {
	//Copy() tokens.Scope
	Eval(key string, args []tokens.Token, ctx context.Context) (tokens.Token, error)
  Permissive() bool
	SetVar(name string, v Var)
}

// must be registered by directives package
var NewLambdaScope func(fnScope tokens.Scope, callerScope tokens.Scope) LambdaScope
