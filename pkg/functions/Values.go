package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func stringDictValues(arg tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertStringDict(arg)
	if err != nil {
		return nil, err
	}

	values := make([]tokens.Token, 0)

	if err := d.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
		values = append(values, value)
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(values, ctx), nil
}

func intDictValues(arg tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertIntDict(arg)
	if err != nil {
		return nil, err
	}

	values := make([]tokens.Token, 0)

	if err := d.Loop(func(key *tokens.Int, value tokens.Token, last bool) error {
		values = append(values, value)
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(values, ctx), nil
}

func Values(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	switch {
	case tokens.IsStringDict(args[0]):
		return stringDictValues(args[0], ctx)
	case tokens.IsIntDict(args[0]):
		return intDictValues(args[0], ctx)
	default:
		errCtx := ctx
		err := errCtx.NewError("Error: expected dict")
		return nil, err
	}
}
