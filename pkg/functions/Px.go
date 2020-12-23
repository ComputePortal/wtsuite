package functions

import (
	"math"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Px(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 arg")
	}

	f, err := tokens.AssertAnyIntOrFloat(args[0])
	if err != nil {
		return nil, err
	}

	switch f.Unit() {
	case "px":
		return f, nil
	case "rem":
		if tokens.PX_PER_REM <= 0.0 {
			return nil, ctx.NewError("Error: rem-per-px not yet set by config.json")
		}

		return tokens.NewValueUnitFloat(math.Round(f.Value()*float64(tokens.PX_PER_REM)), "px", ctx), nil
	default:
		return nil, ctx.NewError("Error: don't know how to convert " + f.Unit() + " to px")
	}
}
