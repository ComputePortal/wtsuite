package functions

import (
	"path/filepath"

	"../tokens/context"
	tokens "../tokens/html"
)

func Dir(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: exepected 1 argument")
	}

	path, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(path.Value())

	return tokens.NewValueString(dir, ctx), nil
}
