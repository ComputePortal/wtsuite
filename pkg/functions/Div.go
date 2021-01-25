package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func divInts(a *tokens.Int, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewValueFloat(float64(a.Value())/float64(b.Value()), ctx), nil
}

func divIntFloat(a *tokens.Int, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if b.Unit() != "" {
		return nil, ctx.NewError("Error: can't divide by unit")
	}

	return tokens.NewValueFloat(float64(a.Value())/b.Value(), ctx), nil
}

func divFloatInt(a *tokens.Float, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewValueUnitFloat(a.Value()/float64(b.Value()), a.Unit(), ctx), nil
}

func divFloats(a *tokens.Float, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	unit := a.Unit()
	if b.Unit() != "" {
		if a.Unit() == b.Unit() {
			unit = ""
		} else {
			return nil, ctx.NewError("Error: can't divide two different units")
		}
	}

	return tokens.NewValueUnitFloat(a.Value()/b.Value(), unit, ctx), nil
}

func Div(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, NewBinaryInterface(ctx))
  if err != nil {
    return nil, err
  }

	switch a := args[0].(type) {
	case *tokens.Int:
		switch b := args[1].(type) {
		case *tokens.Int:
			return divInts(a, b, ctx)
		case *tokens.Float:
			return divIntFloat(a, b, ctx)
		default:
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected Int or Float rhs")
		}
	case *tokens.Float:
		switch b := args[1].(type) {
		case *tokens.Int:
			return divFloatInt(a, b, ctx)
		case *tokens.Float:
			return divFloats(a, b, ctx)
		default:
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected Int or Float rhs")
		}
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected Int or Float lhs")
	}
}
