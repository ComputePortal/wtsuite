package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Not(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	b, err := tokens.AssertBool(args[0])
	if err != nil {
		return nil, err
	}

	return tokens.NewValueBool(!b.Value(), ctx), nil
}
