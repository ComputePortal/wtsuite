package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func filterList(scope tokens.Scope, arg0 tokens.Token, arg1 tokens.Token, ctx context.Context) (tokens.Token, error) {
	list, err := tokens.AssertList(arg0)
	if err != nil {
		context.AppendContextString(err, "Info: used here", ctx)
		return nil, err
	}

	fn, err := AssertFun(arg1)
	if err != nil {
		return nil, err
	}

	innerArgs := make([]tokens.Token, fn.Len())
	result := make([]tokens.Token, 0)

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
		cond, err := fn.EvalFun(scope, innerArgs, ctx)
		if err != nil {
			return err
		}

		b, err := tokens.AssertBool(cond)
		if err != nil {
			return err
		}

		if b.Value() {
			result = append(result, value)
		}
		return err
	}); err != nil {
		return nil, err
	}

	return tokens.NewValuesList(result, ctx), nil
}

// filter keys that are in the list
func filterDict(scope tokens.Scope, arg0 tokens.Token, arg1 tokens.Token, ctx context.Context) (tokens.Token, error) {
	d, err := tokens.AssertStringDict(arg0)
	if err != nil {
		panic(err)
	}

	l_, err := tokens.AssertList(arg1)
	if err != nil {
		return nil, err
	}

	l, err := l_.EvalList(scope)
	if err != nil {
		return nil, err
	}

	res := tokens.NewEmptyStringDict(ctx)

	hasKey := func(test *tokens.String) (bool, error) {
		ok := false
		if err := l.Loop(func(i int, val_ tokens.Token, last bool) error {
			val, err := tokens.AssertString(val_)
			if err != nil {
				return err
			}

			if val.Value() == test.Value() {
				ok = true
			}

			return nil
		}); err != nil {
			return false, err
		}

		return ok, nil
	}

	if err := d.Loop(func(key *tokens.String, val tokens.Token, last bool) error {
		ok, err := hasKey(key)

		if err != nil {
			return err
		}

		if ok {
			res.Set(key, val)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

func Filter(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	arg0, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if tokens.IsStringDict(arg0) {
		return filterDict(scope, arg0, args[1], ctx)
	} else if tokens.IsList(arg0) {
		return filterList(scope, arg0, args[1], ctx)
	} else {
		return nil, ctx.NewError("Error: expected list or dict for first arg")
	}
}
