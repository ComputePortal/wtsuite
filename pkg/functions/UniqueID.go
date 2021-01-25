package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/styles"
)

func UniqueID(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	id := styles.NewUniqueID()

	return tokens.NewString(id, ctx)
}
