package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
	"../tree/styles"
)

func UniqueID(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	id := styles.NewUniqueID()

	return tokens.NewString(id, ctx)
}
