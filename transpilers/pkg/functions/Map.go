package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func Map(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	arg0, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	list, err := tokens.AssertList(arg0)
	if err != nil {
		return nil, err
	}

	arg1, err := args[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	fn, err := AssertFun(arg1)
	if err != nil {
		return nil, err
	}

	if fn.Len() < 0 {
		return nil, ctx.NewError("Error: varargs functions cant be used for map(), (hint: wrap in a new function(){})")
	}

	innerArgs := make([]tokens.Token, fn.Len())
	result := make([]tokens.Token, list.Len())

	if err := list.Loop(func(i int, value tokens.Token, last bool) error {
		switch fn.Len() {
		case 1:
			innerArgs[0] = value
		case 2:
			innerArgs[0] = tokens.NewValueInt(i, ctx)
			innerArgs[1] = value
		default:
			errCtx := fn.Context()
			return errCtx.NewError("Error: unsupported number of fn-args (expected 1 or 2)")
		}

		var err error
		result[i], err = fn.EvalFun(scope, innerArgs, ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(result, ctx), nil
}
