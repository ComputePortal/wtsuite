package functions

import (
	"time"

	"../tokens/context"
	tokens "../tokens/html"
)

func Year(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: unexpected arguments")
	}

	return tokens.NewValueInt(time.Now().Year(), ctx), nil
}
