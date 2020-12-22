package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func addInts(a *tokens.Int, b *tokens.Int, ctx context.Context) (tokens.Token, error) {
	return tokens.NewInt(a.Value()+b.Value(), ctx)
}

func addIntFloat(a *tokens.Int, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if b.Unit() != "" {
		return nil, ctx.NewError("Error: can't add unit to non-unit")
	}

	return tokens.NewValueFloat(float64(a.Value())+b.Value(), ctx), nil
}

func addFloats(a *tokens.Float, b *tokens.Float, ctx context.Context) (tokens.Token, error) {
	if a.Unit() != b.Unit() {
		if tokens.PX_PER_REM > 0.0 {
			if a.Unit() == "px" && b.Unit() == "rem" {
				return tokens.NewValueUnitFloat(a.Value()+b.Value()*float64(tokens.PX_PER_REM),
					"px", ctx), nil
			} else if a.Unit() == "rem" && b.Unit() == "px" {
				return tokens.NewValueUnitFloat(b.Value()+a.Value()*float64(tokens.PX_PER_REM),
					"px", ctx), nil
			}
		}

		return nil, ctx.NewError("Error: units differ")
	}

	return tokens.NewValueUnitFloat(a.Value()+b.Value(), a.Unit(), ctx), nil
}

func joinLists(a *tokens.List, b *tokens.List, ctx context.Context) (tokens.Token, error) {
	res := make([]tokens.Token, 0)

	appendFn := func(i int, v tokens.Token, last bool) error {
		res = append(res, v)
		return nil
	}

	if err := a.Loop(appendFn); err != nil {
		panic(err)
	}

	if err := b.Loop(appendFn); err != nil {
		panic(err)
	}

	return tokens.NewValuesList(res, ctx), nil
}

func copyValue(v_ tokens.Token, ctx context.Context) (tokens.Token, error) {
  switch v := v_.(type) {
  case *tokens.Int:
    return tokens.NewValueInt(v.Value(), ctx), nil
  case *tokens.Float:
    return tokens.NewValueUnitFloat(v.Value(), v.Unit(), ctx), nil
  case *tokens.String:
    return tokens.NewValueString(v.Value(), ctx), nil
  case *tokens.List:
    return v.Copy(ctx)
  case *tokens.StringDict:
    return v.Copy(ctx)
  case *tokens.IntDict:
    return v.Copy(ctx)
  default:
    return nil, ctx.NewError("Error: not addable")
  }
}

func Add(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		err := ctx.NewError("Error: expected 2 arguments")
		return nil, err
	}

  if scope.Permissive() {
    if tokens.IsNull(args[0]) && tokens.IsNull(args[1]) {
      return tokens.NewNull(ctx), nil
    } else if tokens.IsNull(args[0]) {
      return copyValue(args[1], ctx)
    } else if tokens.IsNull(args[1]) {
      return copyValue(args[0], ctx)
    }
  }

	switch a := args[0].(type) {
	case *tokens.Int:
		switch b := args[1].(type) {
		case *tokens.Int:
			return addInts(a, b, ctx)
		case *tokens.Float:
			return addIntFloat(a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected int or float")
		}
	case *tokens.Float:
		switch b := args[1].(type) {
		case *tokens.Int:
			return addIntFloat(b, a, ctx)
		case *tokens.Float:
			return addFloats(a, b, ctx)
		default:
			errCtx := b.Context()
			return nil, errCtx.NewError("Error: expected int or float")
		}
	case *tokens.String:
		switch b := args[1].(type) {
		case *tokens.String:
			return joinStrings(a, b, ctx)
		default:
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected two strings")
		}
	case *tokens.List:
		switch b := args[1].(type) {
		case *tokens.List:
			return joinLists(a, b, ctx)
		default:
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected two lists")
		}
	case *tokens.StringDict:
		switch b := args[1].(type) {
		case *tokens.StringDict:
			return mergeStringDicts(scope, a, b, ctx)
		default:
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected two strings dicts")
		}
	case *tokens.IntDict:
		switch b := args[1].(type) {
		case *tokens.IntDict:
			return mergeIntDicts(scope, a, b, ctx)
		default:
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected two int dicts")
		}
	default:
		errCtx := a.Context()
		return nil, errCtx.NewError("Error: expected int or float")
	}
}
