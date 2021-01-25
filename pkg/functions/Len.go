package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Len(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

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
		res = a.Args().Len()
	case *AnonFun:
		res = a.Len()
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected string, list, dict or function")
	}

	return tokens.NewInt(res, ctx)
}
