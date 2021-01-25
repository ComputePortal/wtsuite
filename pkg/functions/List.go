package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func List(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }
  
	switch len(args) {
	case 1:
		str, err := tokens.AssertString(args[0])
		if err != nil {
			return nil, err
		}

		s := str.Value()
		n := len(s)
		content := make([]tokens.Token, n)

		for i := 0; i < n; i++ {
			content[i] = tokens.NewValueString(s[i:i+1], ctx)
		}

		return tokens.NewValuesList(content, ctx), nil
	case 2:
		n, err := tokens.AssertInt(args[0])
		if err != nil {
			return nil, err
		}

		content := make([]tokens.Token, n.Value())

		for i := 0; i < n.Value(); i++ {
			content[i] = args[1]
		}

		return tokens.NewValuesList(content, ctx), nil
	default:
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}

}
