package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Invert(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	color, err := tokens.AssertColor(args[0])
	if err != nil {
		return nil, err
	}

	r, g, b, a := color.Values()

	ri := 255 - r
	gi := 255 - g
	bi := 255 - b

	return tokens.NewValueColor(ri, gi, bi, a, ctx), nil
}
