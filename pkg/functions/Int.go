package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func castIntToInt(a *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewInt(a.Value(), ctx)
}

func castFloatToInt(a *tokens.Float, ctx context.Context) (tokens.Token, error) {
	// unit is lost
	return tokens.NewInt(int(a.Value()), ctx)
}

func Int(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		ctx.NewError("Error: expected 1 argument")
	}

	switch a := args[0].(type) {
	case *tokens.Int:
		return castIntToInt(a, ctx)
	case *tokens.Float:
		return castFloatToInt(a, ctx)
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected Int or Float")
	}
}
