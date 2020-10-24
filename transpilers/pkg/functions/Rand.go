package functions

import (
	"math/rand"

	"../tokens/context"
	tokens "../tokens/html"
)

func randInt(r *rand.Rand, a *tokens.Int, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	if b.Value() <= a.Value() {
		return nil, ctx.NewError("Error: second value must be greater than first value")
	}

	value := r.Intn(b.Value()-a.Value()) + a.Value()
	return tokens.NewInt(value, ctx)
}

func randFloat(r *rand.Rand, a *tokens.Float, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if a.Unit() != b.Unit() {
		return nil, ctx.NewError("Error: units differ")
	}

	if b.Value() <= a.Value() {
		return nil, ctx.NewError("Error: second value must be greater than first value")
	}

	value := r.Float64()*(b.Value()-a.Value()) + a.Value()
	return tokens.NewValueUnitFloat(value, a.Unit(), ctx), nil
}

func Rand(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 3 {
		ctx.NewError("Error: expected 3 arguments")
	}

	seed, ok := args[0].(*tokens.Int)
	if !ok {
		return nil, ctx.NewError("Error: expected integer seed")
	}

	r := rand.New(rand.NewSource(int64(seed.Value())))

	switch a := args[1].(type) {
	case *tokens.Int:
		switch b := args[2].(type) {
		case *tokens.Int:
			return randInt(r, a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected Int")
		}
	case *tokens.Float:
		switch b := args[2].(type) {
		case *tokens.Float:
			return randFloat(r, a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected Float")
		}
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected Int or Float")
	}
}
