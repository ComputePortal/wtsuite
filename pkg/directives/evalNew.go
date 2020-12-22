package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

func evalNew(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	args, err := functions.EvalArgs(scope, args)
	if err != nil {
		return nil, err
	}

	nameToken, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	if err := AssertValidVar(nameToken); err != nil {
		return nil, err
	}

	valueToken := args[1]

	key := nameToken.Value()
	switch {
	case HasGlobal(key):
		errCtx := nameToken.InnerContext()
		return nil, errCtx.NewError("Error: can't redefine global")
	case scope.HasVar(key):
		errCtx := nameToken.InnerContext()
		return nil, errCtx.NewError("Error: can't redefine variable")
	default:
		v := functions.Var{valueToken, false, false, false, false, ctx}
		scope.SetVar(key, v)
	}

	return valueToken, nil
}
