package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func Not(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	b, err := tokens.AssertBool(args[0])
	if err != nil {
		return nil, err
	}

	return tokens.NewValueBool(!b.Value(), ctx), nil
}
