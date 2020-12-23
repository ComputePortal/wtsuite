package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Xor(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected two arguments")
	}

	a, err := tokens.AssertBool(args[0])
	if err != nil {
		return nil, err
	}

	b, err := tokens.AssertBool(args[1])
	if err != nil {
		return nil, err
	}

	res := a.Value() != b.Value()

	return tokens.NewBool(res, ctx)
}
