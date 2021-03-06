package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func castIntToFloat(a *tokens.Int, unit string, ctx context.Context) (tokens.Token, error) {
	return tokens.NewValueUnitFloat(float64(a.Value()), unit, ctx), nil
}

func castFloatToFloat(a *tokens.Float, unit string, ctx context.Context) (tokens.Token, error) {
	return tokens.NewValueUnitFloat(a.Value(), unit, ctx), nil
}

func Float(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }
  
	unit := ""
	if len(args) == 2 {
		u, ok := args[1].(*tokens.String)
		if !ok {
			errCtx := args[1].Context()
			return nil, errCtx.NewError("Error: expected string")
		}
		unit = u.Value()
	} else if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}

	switch a := args[0].(type) {
	case *tokens.Int:
		return castIntToFloat(a, unit, ctx)
	case *tokens.Float:
		return castFloatToFloat(a, unit, ctx)
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected Int or Float")
	}
}
