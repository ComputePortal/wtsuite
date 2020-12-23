package functions

import (
	"math"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

// a%b -> remainder
func Mod(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	a, err := tokens.AssertInt(args[0])
	if err != nil {
		return nil, err
	}

	b, err := tokens.AssertInt(args[1])
	if err != nil {
		return nil, err
	}

	cVal := math.Mod(float64(a.Value()), float64(b.Value()))

	return tokens.NewValueInt(int(cVal), ctx), nil
}
