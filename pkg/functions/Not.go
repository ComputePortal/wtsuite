package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Not(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, NewUnaryInterface(ctx))
  if err != nil {
    return nil, err
  }

	b, err := tokens.AssertBool(args[0])
	if err != nil {
		return nil, err
	}

	return tokens.NewValueBool(!b.Value(), ctx), nil
}
