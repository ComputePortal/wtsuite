package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

// - convert a key list and value list into a dict
func stringKeysValuesToDict(a *tokens.List, b *tokens.List, ctx context.Context) (tokens.Token, error) {
	result := tokens.NewEmptyStringDict(ctx)

	if err := a.Loop(func(i int, key tokens.Token, last bool) error {
		result.Set(key, b.GetTokens()[i])
		return nil
	}); err != nil {
		panic(err)
	}

	return result, nil
}

func intKeysValuesToDict(a *tokens.List, b *tokens.List, ctx context.Context) (tokens.Token, error) {
	result := tokens.NewEmptyIntDict(ctx)

	if err := a.Loop(func(i int, key tokens.Token, last bool) error {
		result.Set(key, b.GetTokens()[i])
		return nil
	}); err != nil {
		panic(err)
	}

	return result, nil
}

func Dict(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	switch len(args) {
	case 1:
		pairs, err := tokens.AssertList(args[0])
		if err != nil {
			return nil, err
		}

		keys := tokens.NewEmptyList(ctx)
		values := tokens.NewEmptyList(ctx)

		isInts := false
		isStrings := false
		if err := pairs.Loop(func(i int, v tokens.Token, last bool) error {
			pair, err := tokens.AssertList(v)
			if err != nil {
				return err
			}

			if pair.Len() != 2 {
				errCtx := ctx
				return errCtx.NewError("Error: expected pair")
			}

			ab := pair.GetTokens()

			if tokens.IsString(ab[0]) {
				if isInts {
					return ctx.NewError("Error: expected all int or all string keys")
				}

				isStrings = true
			} else if tokens.IsInt(ab[0]) {
				if isStrings {
					return ctx.NewError("Error: expected all int or all string keys")
				}

				isInts = true
			} else {
				return ctx.NewError("Error: expected int or string keys")
			}
			keys.Append(ab[0])
			values.Append(ab[1])
			return nil
		}); err != nil {
			return nil, err
		}

		if isInts {
			return intKeysValuesToDict(keys, values, ctx)
		} else if isStrings {
			return stringKeysValuesToDict(keys, values, ctx)
		} else {
			return tokens.NewEmptyStringDict(ctx), nil
		}
	case 2:
		a, err := tokens.AssertList(args[0])
		if err != nil {
			return nil, err
		}

		b, err := tokens.AssertList(args[1])
		if err != nil {
			return nil, err
		}

		if a.Len() != b.Len() {
			errCtx := ctx
			return nil, errCtx.NewError("Error: key list and value list don't have same length")
		}

		switch {
		case tokens.IsStringList(a):
			return stringKeysValuesToDict(a, b, ctx)
		case tokens.IsIntList(a):
			return intKeysValuesToDict(a, b, ctx)
		default:
			errCtx := a.Context()
			err := errCtx.NewError("Error: expected Int list or String list")
			err.AppendContextString("Info: needed here", ctx)
			return nil, err
		}
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 2 arguments")
	}
}
