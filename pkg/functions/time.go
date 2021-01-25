package functions

import (
	"time"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Year(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args := args_.Values()

	if len(args) != 0 {
		return nil, ctx.NewError("Error: unexpected arguments")
	}

	return tokens.NewValueInt(time.Now().Year(), ctx), nil
}
