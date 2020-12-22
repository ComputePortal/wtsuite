package functions

import (
	"../files"
	"../tokens/context"
	tokens "../tokens/html"
)

func Path(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	s, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	path, err := files.Search(ctx.Path(), s.Value())
	if err != nil {
		return nil, ctx.NewError("Error: couldn't find file " + s.Value())
	}

	return tokens.NewValueString(path, s.Context()), nil
}
