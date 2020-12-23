package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func negInt(a *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewInt(-a.Value(), ctx)
}

func negFloat(a *tokens.Float, ctx context.Context) (tokens.Token, error) {
	return tokens.NewValueUnitFloat(-a.Value(), a.Unit(), ctx), nil
}

func Neg(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	switch a := args[0].(type) {
	case *tokens.Int:
		return negInt(a, ctx)
	case *tokens.Float:
		return negFloat(a, ctx)
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected Int or Float")
	}
}
