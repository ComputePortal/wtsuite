package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func And(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, NewBinaryInterface(ctx))
  if err != nil {
    return nil, err
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
