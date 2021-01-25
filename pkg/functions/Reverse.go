package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Reverse(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

  if len(args) != 1 {
    return nil, ctx.NewError("Error: expected 1 argument")
  }

  lst, err := tokens.AssertList(args[0])
  if err != nil {
    return nil, err
  }

  n := lst.Len()
  vals := make([]tokens.Token, n)

  if err := lst.Loop(func(i int, value tokens.Token, last bool) error {
    vals[n-1-i] = value
    return nil
  }); err != nil {
    panic(err)
  }

  return tokens.NewValuesList(vals, ctx), nil
}
