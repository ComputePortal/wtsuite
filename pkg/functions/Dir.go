package functions

import (
	"path/filepath"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Dir(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	path, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(path.Value())

	return tokens.NewValueString(dir, ctx), nil
}
