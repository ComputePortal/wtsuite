package functions

import (
	"math"

	"../tokens/context"
	tokens "../tokens/html"
)

func subColorBrightness(color *tokens.Color, d float64,
	ctx context.Context) (tokens.Token, error) {
	// alpha remains unchanged

	r, g, b, a := color.Values()

	r = int(math.Max(0, float64(r)-math.Round(d)))
	g = int(math.Max(0, float64(g)-math.Round(d)))
	b = int(math.Max(0, float64(b)-math.Round(d)))

	return tokens.NewColor(r, g, b, a, ctx)
}

func Darken(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	color, ok := args[0].(*tokens.Color)
	if !ok {
		errCtx := args[0].Context()
		return nil, errCtx.NewError("Error: expected color")
	}

	factor, ok := args[1].(*tokens.Float)
	if !ok {
		errCtx := args[1].Context()
		return nil, errCtx.NewError("Error: expected float")
	}

	switch factor.Unit() {
	case "":
		return subColorBrightness(color, factor.Value()*256.0, ctx)
	case "%":
		return subColorBrightness(color, factor.Value()*256.0/100.0, ctx)
	default:
		errCtx := args[1].Context()
		return nil, errCtx.NewError("Error: expected unitless or '%' float")
	}
}
