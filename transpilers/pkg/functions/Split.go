package functions

import (
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
)

func Split(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	str, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	sep, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	res := strings.Split(str.Value(), sep.Value())

	lst := tokens.NewValuesList(res, ctx)

	return lst, nil
}
