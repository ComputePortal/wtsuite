package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func stringDictItems(arg tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertStringDict(arg)
	if err != nil {
		return nil, err
	}

	items := make([]tokens.Token, 0)

	if err := d.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
		item := tokens.NewValuesList([]tokens.Token{key, value}, key.Context())
		items = append(items, item)
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(items, ctx), nil
}

func intDictItems(arg tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertIntDict(arg)
	if err != nil {
		return nil, err
	}

	items := make([]tokens.Token, 0)

	if err := d.Loop(func(key *tokens.Int, value tokens.Token, last bool) error {
		item := tokens.NewValuesList([]tokens.Token{key, value}, key.Context())
		items = append(items, item)
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(items, ctx), nil
}

func Items(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	switch {
	case tokens.IsStringDict(args[0]):
		return stringDictItems(args[0], ctx)
	case tokens.IsIntDict(args[0]):
		return intDictItems(args[0], ctx)
	default:
		errCtx := ctx
		err := errCtx.NewError("Error: expected dict")
		return nil, err
	}
}
