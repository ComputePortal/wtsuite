package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func And(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	a_, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	a, err := tokens.AssertBool(a_)
	if err != nil {
		return nil, err
	}

	if !a.Value() {
		// shortcircuit evaluation
		return tokens.NewBool(false, ctx)
	}

	b_, err := args[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	b, err := tokens.AssertBool(b_)
	if err != nil {
		return nil, err
	}

	return tokens.NewBool(b.Value(), ctx)
}
