package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func Len(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	res := 0
	switch a := args[0].(type) {
	case *tokens.StringDict:
		res = a.Len()
	case *tokens.IntDict:
		res = a.Len()
	case *tokens.RawDict:
		res = a.Len()
	case *tokens.List:
		res = a.Len()
	case *tokens.String:
		res = len(a.Value())
	case *tokens.Function:
		res = len(a.Args())
	case *AnonFun:
		res = a.Len()
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected string, list, dict or function")
	}

	return tokens.NewInt(res, ctx)
}
