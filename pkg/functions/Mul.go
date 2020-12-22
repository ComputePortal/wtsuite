package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func mulInts(a *tokens.Int, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewInt(a.Value()*b.Value(), ctx)
}

func mulIntFloat(a *tokens.Int, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	return tokens.NewValueUnitFloat(float64(a.Value())*b.Value(), b.Unit(), ctx), nil
}

func mulFloats(a *tokens.Float, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if a.Unit() != "" && b.Unit() != "" {
		return nil, ctx.NewError("Error: can't multiply 2 units")
	}
	value := a.Value() * b.Value()
	unit := ""
	if a.Unit() != "" {
		unit = a.Unit()
	} else if b.Unit() != "" {
		unit = b.Unit()
	}

	return tokens.NewValueUnitFloat(value, unit, ctx), nil
}

func Mul(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	switch a := args[0].(type) {
	case *tokens.Int:
		switch b := args[1].(type) {
		case *tokens.Int:
			return mulInts(a, b, ctx)
		case *tokens.Float:
			return mulIntFloat(a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected Int or Float")
		}
	case *tokens.Float:
		switch b := args[1].(type) {
		case *tokens.Int:
			return mulIntFloat(b, a, ctx)
		case *tokens.Float:
			return mulFloats(a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected Int or Float")
		}
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected Int or Float")
	}
}
