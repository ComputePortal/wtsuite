package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func intsToListToken(ints []int, ctx context.Context) *tokens.List {
	list := make([]tokens.Token, len(ints))

	for i, v := range ints {
		list[i] = tokens.NewValueInt(v, ctx)
	}

	return tokens.NewValuesList(list, ctx)
}

func floatsToListToken(floats []float64, unit string, ctx context.Context) *tokens.List {
	list := make([]tokens.Token, len(floats))

	for i, v := range floats {
		list[i] = tokens.NewValueUnitFloat(v, unit, ctx)
	}

	return tokens.NewValuesList(list, ctx)
}

func seqLastInt(last *tokens.Int, ctx context.Context) (tokens.Token, error) {
	res := make([]int, 0)

	for i := 0; i < last.Value(); i++ {
		res = append(res, i)
	}

	return intsToListToken(res, ctx), nil
}

func seqLastFloat(last *tokens.Float, ctx context.Context) (tokens.Token, error) {
	res := make([]float64, 0)

	for i := 0; float64(i) < last.Value(); i++ {
		res = append(res, float64(i))
	}

	return floatsToListToken(res, last.Unit(), ctx), nil
}

func seqFirstLastInt(first *tokens.Int, last *tokens.Int, ctx context.Context) (tokens.Token, error) {
	res := make([]int, 0)

	for i := first.Value(); i < last.Value(); i++ {
		res = append(res, i)
	}

	return intsToListToken(res, ctx), nil
}

func seqFirstLastFloat(first *tokens.Float, last *tokens.Float, ctx context.Context) (tokens.Token, error) {
	res := make([]float64, 0)

	for i := first.Value(); i < last.Value(); i += 1.0 {
		res = append(res, i)
	}

	return floatsToListToken(res, first.Unit(), ctx), nil
}

func seqFirstIncrLastInt(first *tokens.Int, incr *tokens.Int, last *tokens.Int, ctx context.Context) (tokens.Token, error) {
	res := make([]int, 0)

	if incr.Value() == 0 {
		errCtx := context.MergeContexts(incr.Context(), ctx)
		return nil, errCtx.NewError("Error: incr can't be 0")
	} else if incr.Value() > 0 {
		for i := first.Value(); i < last.Value(); i += incr.Value() {
			res = append(res, i)
		}
	} else {
		d := -incr.Value()
		for i := first.Value(); i > last.Value(); i -= d {
			res = append(res, i)
		}
	}

	return intsToListToken(res, ctx), nil
}

func seqFirstIncrLastFloat(first *tokens.Float, incr *tokens.Float, last *tokens.Float, ctx context.Context) (tokens.Token, error) {
	res := make([]float64, 0)

	if incr.Value() == 0.0 {
		errCtx := context.MergeContexts(incr.Context(), ctx)
		return nil, errCtx.NewError("Error: incr can't be 0")
	} else if incr.Value() > 0.0 {
		for i := first.Value(); i < last.Value(); i += incr.Value() {
			res = append(res, i)
		}
	} else {
		d := -incr.Value()
		for i := first.Value(); i > last.Value(); i -= d {
			res = append(res, i)
		}
	}

	return floatsToListToken(res, first.Unit(), ctx), nil
}

// seq(last) // starting at 0 and going up
// seq(first, last) // increasing
// seq(first, incr, last) (similar to [first:incr:last] in formula form)
func Seq(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	intArgs := make([]*tokens.Int, len(args))
	for i, arg := range args {
		if tokens.IsFloat(arg) {
			return seqFloat(args, ctx)
		}
		var err error
		intArgs[i], err = tokens.AssertInt(arg)
		if err != nil {
			return nil, err
		}
	}

	switch len(intArgs) {
	case 1:
		return seqLastInt(intArgs[0], ctx)
	case 2:
		return seqFirstLastInt(intArgs[0], intArgs[1], ctx)
	case 3:
		return seqFirstIncrLastInt(intArgs[0], intArgs[1], intArgs[2], ctx)
	default:
		return nil, ctx.NewError("Error: expected 1, 2 or 3 arguments")
	}
}

func seqFloat(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	floatArgs := make([]*tokens.Float, len(args))
	for i, arg := range args {
		var err error
		floatArgs[i], err = tokens.AssertAnyIntOrFloat(arg)
		if err != nil {
			return nil, err
		}
	}

	if len(floatArgs) < 1 {
		return nil, ctx.NewError("Error: expected 1, 2 or 3 arguments")
	}

	// check that all units are the same
	firstUnit := floatArgs[0].Unit()
	for i, f := range floatArgs {
		if i > 0 {
			if f.Unit() != firstUnit {
				return nil, ctx.NewError("Error: differnt arg float units, first got " + firstUnit + ", then got " + f.Unit())
			}
		}
	}

	switch len(floatArgs) {
	case 1:
		return seqLastFloat(floatArgs[0], ctx)
	case 2:
		return seqFirstLastFloat(floatArgs[0], floatArgs[1], ctx)
	case 3:
		return seqFirstIncrLastFloat(floatArgs[0], floatArgs[1], floatArgs[2], ctx)
	default:
		return nil, ctx.NewError("Error: expected 1, 2 or 3 arguments")
	}
}
