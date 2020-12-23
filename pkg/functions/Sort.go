package functions

import (
	"sort"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func sortInts(list *tokens.List, ctx context.Context) (*tokens.List, error) {
	values := make([]*tokens.Int, list.Len())

	if err := list.Loop(func(i int, val tokens.Token, last bool) error {
		value, err := tokens.AssertInt(val)
		if err != nil {
			return err
		}

		values[i] = value
		return nil
	}); err != nil {
		return nil, err
	}

	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Value() < values[j].Value()
	})

	return tokens.NewValuesList(values, ctx), nil
}

func sortFloats(first *tokens.Float, list *tokens.List, ctx context.Context) (*tokens.List, error) {
	values := make([]*tokens.Float, list.Len())

	if err := list.Loop(func(i int, val tokens.Token, last bool) error {
		value, err := tokens.AssertFloat(val, first.Unit())
		if err != nil {
			return err
		}

		values[i] = value
		return nil
	}); err != nil {
		return nil, err
	}

	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Value() < values[j].Value()
	})

	return tokens.NewValuesList(values, ctx), nil
}

func sortStrings(list *tokens.List, ctx context.Context) (*tokens.List, error) {
	values := make([]*tokens.String, list.Len())

	if err := list.Loop(func(i int, val tokens.Token, last bool) error {
		value, err := tokens.AssertString(val)
		if err != nil {
			return err
		}

		values[i] = value
		return nil
	}); err != nil {
		return nil, err
	}

	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Value() < values[j].Value()
	})

	return tokens.NewValuesList(values, ctx), nil
}

func sortGeneric(scope tokens.Scope, list *tokens.List, fn Fun, ctx context.Context) (*tokens.List, error) {
	values := make([]tokens.Token, list.Len())

	if err := list.Loop(func(i int, val tokens.Token, last bool) error {
		values[i] = val
		return nil
	}); err != nil {
		panic(err)
	}

	innerArgs := make([]tokens.Token, fn.Len())

	var sortErr error = nil
	sort.SliceStable(values, func(i, j int) bool {
		if sortErr != nil {
			return true // quit as quickly as possible
		}

		switch fn.Len() {
		case 1:
			innerArgs[0] = tokens.NewValueInt(i, ctx)
			a, err := fn.EvalFun(scope, innerArgs, ctx)
			if err != nil {
				sortErr = err
				return true
			}

			innerArgs[1] = tokens.NewValueInt(j, ctx)
			b, err := fn.EvalFun(scope, innerArgs, ctx)
			if err != nil {
				sortErr = err
				return true
			}

			cond, err := LT(scope, []tokens.Token{a, b}, ctx)
			if err != nil {
				sortErr = err
				return true
			}

			cond_, err := tokens.AssertBool(cond)
			if err != nil {
				panic(err)
			}

			return cond_.Value()
		case 2:
			innerArgs[0] = values[i]
			innerArgs[1] = values[j]
			cond, err := fn.EvalFun(scope, innerArgs, ctx)
			if err != nil {
				sortErr = err
				return true
			} else {
				b, err := tokens.AssertBool(cond)
				if err != nil {
					sortErr = err
					return true
				} else {
					return b.Value()
				}
			}
		default:
			errCtx := fn.Context()
			sortErr = errCtx.NewError("Error: unsupported number of fn-args (expected 1 or 2)")
			return true
		}
	})

	if sortErr != nil {
		return nil, sortErr
	}

	return tokens.NewValuesList(values, ctx), nil
}

func Sort(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) == 0 {
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}

	arg0, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	list, err := tokens.AssertList(arg0)
	if err != nil {
		return nil, err
	}

	if list.Len() < 2 {
		return list, nil
	}

	first, err := list.Get(0)
	if err != nil {
		panic(err)
	}

	switch len(args) {
	case 1:
		switch f := first.(type) {
		case *tokens.Int:
			return sortInts(list, ctx)
		case *tokens.Float:
			return sortFloats(f, list, ctx)
		case *tokens.String:
			return sortStrings(list, ctx)
		default:
			return nil, ctx.NewError("Error: unsortable types")
		}
	case 2:
		fn, err := AssertFun(args[1])
		if err != nil {
			return nil, err
		}
		return sortGeneric(scope, list, fn, ctx)
	default:
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}
}
