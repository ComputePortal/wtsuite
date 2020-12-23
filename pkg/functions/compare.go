package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func compare(args []tokens.Token, ctx context.Context,
	cii func(int, int) bool,
	cif func(int, float64) bool,
	cfi func(float64, int) bool,
	cff func(float64, float64) bool,
	css func(string, string) bool) (tokens.Token, error) {

	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	res := false
	switch a := args[0].(type) {
	case *tokens.Int:
		switch b := args[1].(type) {
		case *tokens.Int:
			res = cii(a.Value(), b.Value())
		case *tokens.Float:
			if b.Unit() != "" {
				errCtx := b.Context()
				return nil, errCtx.NewError("Error: comparing non-unit to unit")
			}
			res = cif(a.Value(), b.Value())
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected int or float")
		}
	case *tokens.Float:
		switch b := args[1].(type) {
		case *tokens.Int:
			if a.Unit() != "" {
				errCtx := a.Context()
				return nil, errCtx.NewError("Error: comparing unit to non-unit")
			} else {
				res = cfi(a.Value(), b.Value())
			}
		case *tokens.Float:
			if a.Unit() != b.Unit() {
				errCtx := context.MergeContexts(a.Context(), b.Context())
				return nil, errCtx.NewError("Error: units differ")
			}
			res = cff(a.Value(), b.Value())
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected int or float")
		}
	case *tokens.String:
		switch b := args[1].(type) {
		case *tokens.String:
			res = css(a.Value(), b.Value())
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected two strings, or numbers")
		}
	}

	return tokens.NewBool(res, ctx)
}

func LT(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return compare(args, ctx,
		func(a int, b int) bool { return a < b },
		func(a int, b float64) bool { return float64(a) < b },
		func(a float64, b int) bool { return a < float64(b) },
		func(a float64, b float64) bool { return a < b },
		func(a string, b string) bool { return a < b },
	)
}

func LE(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return compare(args, ctx,
		func(a int, b int) bool { return a <= b },
		func(a int, b float64) bool { return float64(a) <= b },
		func(a float64, b int) bool { return a <= float64(b) },
		func(a float64, b float64) bool { return a <= b },
		func(a string, b string) bool { return a <= b },
	)
}

func GT(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return compare(args, ctx,
		func(a int, b int) bool { return a > b },
		func(a int, b float64) bool { return float64(a) > b },
		func(a float64, b int) bool { return a > float64(b) },
		func(a float64, b float64) bool { return a > b },
		func(a string, b string) bool { return a > b },
	)
}

func GE(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return compare(args, ctx,
		func(a int, b int) bool { return a >= b },
		func(a int, b float64) bool { return float64(a) >= b },
		func(a float64, b int) bool { return a >= float64(b) },
		func(a float64, b float64) bool { return a >= b },
		func(a string, b string) bool { return a >= b },
	)
}

func EQ(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
  if len(args) == 2 {
    if tokens.IsNull(args[1]) && !tokens.IsNull(args[0]) {
      return tokens.NewBool(false, ctx)
    } else if tokens.IsNull(args[0]) && !tokens.IsNull(args[1]) {
      return tokens.NewBool(false, ctx)
    } else if tokens.IsNull(args[0]) && tokens.IsNull(args[1]) {
      return tokens.NewBool(true, ctx)
    }
  }

	return compare(args, ctx,
		func(a int, b int) bool { return a == b },
		func(a int, b float64) bool { return float64(a) == b },
		func(a float64, b int) bool { return a == float64(b) },
		func(a float64, b float64) bool { return a == b },
		func(a string, b string) bool { return a == b },
	)
}

func NE(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
  if len(args) == 2 {
    if tokens.IsNull(args[1]) && !tokens.IsNull(args[0]) {
      return tokens.NewBool(true, ctx)
    } else if tokens.IsNull(args[0]) && !tokens.IsNull(args[1]) {
      return tokens.NewBool(true, ctx)
    } else if tokens.IsNull(args[0]) && tokens.IsNull(args[1]) {
      return tokens.NewBool(false, ctx)
    }
  }

	return compare(args, ctx,
		func(a int, b int) bool { return a != b },
		func(a int, b float64) bool { return float64(a) != b },
		func(a float64, b int) bool { return a != float64(b) },
		func(a float64, b float64) bool { return a != b },
		func(a string, b string) bool { return a != b },
	)
}

func minMax(scope tokens.Scope, args []tokens.Token, ctx context.Context, isMax bool) (tokens.Token, error) {
	cond, err := LT(scope, args, ctx)
	if err != nil {
		return nil, err
	}

	b, err := tokens.AssertBool(cond)
	if err != nil {
		panic(err)
	}

	if b.Value() != isMax {
		return args[0], nil
	} else {
		return args[1], nil
	}
}

func Min(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return minMax(scope, args, ctx, false)
}

func Max(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return minMax(scope, args, ctx, true)
}

func IsSame(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	// different from eq because it can be used for any type, and does a deep comparison
	// ints cannot be compared to floats!
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	return tokens.NewValueBool(args[0].IsSame(args[1]), ctx), nil
}
