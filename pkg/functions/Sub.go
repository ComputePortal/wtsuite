package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func subInts(a *tokens.Int, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewInt(a.Value()-b.Value(), ctx)
}

func subIntFloat(a *tokens.Int, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if b.Unit() != "" {
		return nil, ctx.NewError("Error: can't sub unit from non-unit")
	}

	return tokens.NewValueFloat(float64(a.Value())-b.Value(), ctx), nil
}

func subFloatInt(a *tokens.Float, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	if a.Unit() != "" {
		return nil, ctx.NewError("Error: can't sub non-unit from unit")
	}

	return tokens.NewValueFloat(a.Value()-float64(b.Value()), ctx), nil
}

func subFloats(a *tokens.Float, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if a.Unit() != b.Unit() {
		if tokens.PX_PER_REM > 0.0 {
			if a.Unit() == "px" && b.Unit() == "rem" {
				return tokens.NewValueUnitFloat(a.Value()-b.Value()*float64(tokens.PX_PER_REM),
					"px", ctx), nil
			} else if a.Unit() == "rem" && b.Unit() == "px" {
				return tokens.NewValueUnitFloat(a.Value()*float64(tokens.PX_PER_REM)-b.Value(),
					"px", ctx), nil
			}
		}

		return nil, ctx.NewError("Error: units differ")
	}

	return tokens.NewValueUnitFloat(a.Value()-b.Value(), a.Unit(), ctx), nil
}

func Sub(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	switch a := args[0].(type) {
	case *tokens.Int:
		switch b := args[1].(type) {
		case *tokens.Int:
			return subInts(a, b, ctx)
		case *tokens.Float:
			return subIntFloat(a, b, ctx)
		default:
			errCtx := context.MergeContexts(ctx, a.Context(), b.Context())
			return nil, errCtx.NewError("Error: unable to subtract")
		}
	case *tokens.Float:
		switch b := args[1].(type) {
		case *tokens.Int:
			return subFloatInt(a, b, ctx)
		case *tokens.Float:
			return subFloats(a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected Int or Float")
		}
	default:
		errCtx := context.MergeContexts(ctx, args[0].Context(), args[1].Context())
		return nil, errCtx.NewError("Error: unable to subtract")
	}
}
