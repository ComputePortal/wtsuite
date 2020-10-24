package functions

import (
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
)

func doSlice(obj *tokens.List, start, incr, stop *tokens.Int, ctx context.Context) (*tokens.List, error) {
	input := obj.GetTokens()

	result := make([]tokens.Token, 0)

	if incr.Value() < 0 {
		for i := start.Value(); i > stop.Value(); i += incr.Value() {
			result = append(result, input[i])
		}
	} else {
		for i := start.Value(); i < stop.Value(); i += incr.Value() {
			result = append(result, input[i])
		}
	}

	return tokens.NewValuesList(result, ctx), nil
}

// slice(obj, stop) // starting at 0 and going up
// slice(obj, start, stop) // increasing
// slice(obj, (start|null), (incr|null), (stop|null))
// start defaults to 0 if incr > 0 and len(obj)-1 if incr < 0
// incr default to 1, incr==0 is an error
// stop default to len(obj) if incr > 0 and -1 if incr < 0
// negative start or stop is converted to len(obj)-...
func sliceList(args []tokens.Token, ctx context.Context) (*tokens.List, error) {
	if len(args) < 2 || len(args) > 4 {
		return nil, ctx.NewError("Error: expected 2, 3 or 4 arguments")
	}

	obj, err := tokens.AssertList(args[0])
	if err != nil {
		return nil, err
	}

	n := tokens.NewValueInt(obj.Len(), ctx)

	var start *tokens.Int = tokens.NewValueInt(0, ctx)
	var incr *tokens.Int = tokens.NewValueInt(1, ctx)
	var stop *tokens.Int = n

	adjustNegative := func(t tokens.Token) (*tokens.Int, error) {
		i, err := tokens.AssertInt(t)
		if err != nil {
			return nil, err
		}

		if i.Value() < 0 {
			return tokens.NewValueInt(n.Value()+i.Value(), i.Context()), nil
		} else {
			return i, nil
		}
	}

	switch len(args) {
	case 2:
		if !tokens.IsNull(args[1]) {
			stop, err = adjustNegative(args[1])
			if err != nil {
				return nil, err
			}
		}
	case 3:
		if !tokens.IsNull(args[1]) {
			start, err = adjustNegative(args[1])
			if err != nil {
				return nil, err
			}
		}

		if !tokens.IsNull(args[2]) {
			stop, err = adjustNegative(args[2])
			if err != nil {
				return nil, err
			}
		}
	case 4:
		if !tokens.IsNull(args[2]) {
			incr, err = tokens.AssertInt(args[2])
			if err != nil {
				return nil, err
			}

			if incr.Value() == 0 {
				errCtx := incr.Context()
				err := errCtx.NewError("Error: illegal 0 incr")
				err.AppendContextString("Info: used here", ctx)
				return nil, err
			} else if incr.Value() < 0 {
				start = tokens.NewValueInt(n.Value()-1, ctx)
				// stop at -1
				stop = tokens.NewValueInt(-1, ctx)
			}
		}

		if !tokens.IsNull(args[1]) {
			start, err = adjustNegative(args[1])
			if err != nil {
				return nil, err
			}
		}

		if !tokens.IsNull(args[3]) {
			stop, err = adjustNegative(args[3])
			if err != nil {
				return nil, err
			}
		}
	default:
		panic("unhandled")
	}

	if start.Value() > n.Value()-1 || start.Value() < 0 {
		errCtx := start.Context()
		return nil, errCtx.NewError("Error: start index out of range")
	}

	if stop.Value() > n.Value() || stop.Value() < -1 {
		errCtx := start.Context()
		return nil, errCtx.NewError("Error: stop index out of range")
	}

	return doSlice(obj, start, incr, stop, ctx)
}

func sliceString(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	strToken, err := tokens.AssertString(args[0])
	if err != nil {
		panic(err)
	}

	str := strToken.Value()
	dummyList := make([]tokens.Token, len(str))
	for i, _ := range str {
		char := str[i : i+1]

		dummyList[i] = tokens.NewValueString(char, ctx)
	}

	dummyListToken := tokens.NewValuesList(dummyList, ctx)

	lstArgs := append([]tokens.Token{dummyListToken}, args[1:]...)

	lstRes, err := sliceList(lstArgs, ctx)
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	if err := lstRes.Loop(func(i int, v tokens.Token, last bool) error {
		charToken, err := tokens.AssertString(v)
		if err != nil {
			panic(err)
		}

		b.WriteString(charToken.Value())
		return nil
	}); err != nil {
		panic(err)
	}

	return tokens.NewValueString(b.String(), ctx), nil
}

func Slice(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	switch {
	case len(args) == 0:
		return nil, ctx.NewError("Error: expected 2, 3 or 4 arguments")
	case tokens.IsList(args[0]):
		return sliceList(args, ctx)
	case tokens.IsString(args[0]):
		return sliceString(args, ctx)
	default:
		return nil, ctx.NewError("Error: expected list or string")
	}
}
