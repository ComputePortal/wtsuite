package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func stringDictKeys(arg tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertStringDict(arg)
	if err != nil {
		return nil, err
	}

	keys := make([]tokens.Token, 0)

	if err := d.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
		keys = append(keys, key)
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(keys, ctx), nil
}

func intDictKeys(arg tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertIntDict(arg)
	if err != nil {
		return nil, err
	}

	keys := make([]tokens.Token, 0)

	if err := d.Loop(func(key *tokens.Int, value tokens.Token, last bool) error {
		keys = append(keys, key)
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(keys, ctx), nil
}

func Keys(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	switch {
	case tokens.IsStringDict(args[0]):
		return stringDictKeys(args[0], ctx)
	case tokens.IsIntDict(args[0]):
		return intDictKeys(args[0], ctx)
	default:
		errCtx := ctx
		err := errCtx.NewError("Error: expected dict")
		return nil, err
	}
}
