package directives

import (
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
)

func evalAttrEnum(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	comp := func(a, b string) bool { return a == b }
	mod := ""

	switch len(args) {
	case 2:
		return evalThisAttrEnum(scope, args, comp, mod, ctx)
	case 3:
		return evalParentAttrEnum(scope, args, comp, mod, ctx)
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 2 or 3 arguments")
	}
}

func evalAttrEnumS(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	comp := func(a, b string) bool { return strings.HasPrefix(a, b) }
	mod := "^"

	switch len(args) {
	case 2:
		return evalThisAttrEnum(scope, args, comp, mod, ctx)
	case 3:
		return evalParentAttrEnum(scope, args, comp, mod, ctx)
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 2 or 3 arguments")
	}
}

func evalAttrEnumE(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	comp := func(a, b string) bool { return strings.HasSuffix(a, b) }
	mod := "$"

	switch len(args) {
	case 2:
		return evalThisAttrEnum(scope, args, comp, mod, ctx)
	case 3:
		return evalParentAttrEnum(scope, args, comp, mod, ctx)
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 2 or 3 arguments")
	}
}

func evalAttrEnumA(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	switch len(args) {
	case 1:
		return evalThisAttrEnumA(scope, args, ctx)
	case 2:
		return evalParentAttrEnumA(scope, args, ctx)
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 1 or 2 arguments")
	}
}

func evalAttrEnumC(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	comp := func(a, b string) bool { return strings.Contains(a, b) }
	mod := "*"

	switch len(args) {
	case 2:
		return evalThisAttrEnum(scope, args, comp, mod, ctx)
	case 3:
		return evalParentAttrEnum(scope, args, comp, mod, ctx)
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 2 or 3 arguments")
	}
}

func enumMatch(lst *tokens.List, fn func(string) bool) bool {
	found := false

	if err := lst.Loop(func(i int, v tokens.Token, last bool) error {
		if found {
			return nil
		}

		if tokens.IsNull(v) {
			if i != 0 {
				panic("should've been caught before")
			}
			return nil
		}

		str, err := tokens.AssertString(v)
		if err != nil {
			panic(err)
		}

		if fn(str.Value()) {
			found = true
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return found
}

// mod can be "", "*", "^", "$"
func evalThisAttrEnum(scope Scope, args []tokens.Token, fn func(a, b string) bool,
	mod string, ctx context.Context) (tokens.Token, error) {
	arg0, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	arg1, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	lst, err := scope.GetNode().SearchAttrEnum(nil, arg0, ctx)
	if err != nil {
		return nil, err
	}

	found := enumMatch(lst, func(s string) bool { return fn(s, arg1.Value()) })

	if !found {
		return nil, ctx.NewError("Error: attr enum doesn't have match for " + arg1.Value())
	}

	return tokens.NewValueString("@filter "+arg0.Value()+mod+"='"+arg1.Value()+"'", ctx), nil
}

func evalParentAttrEnum(scope Scope, args []tokens.Token, fn func(a, b string) bool,
	mod string, ctx context.Context) (tokens.Token, error) {
	arg0, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	arg1, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	arg2, err := tokens.AssertString(args[2])

	lst, err := scope.GetNode().SearchAttrEnum(arg0, arg1, ctx)
	if err != nil {
		return nil, err
	}

	found := enumMatch(lst, func(s string) bool { return fn(s, arg2.Value()) })

	if !found {
		return nil, ctx.NewError("Error: attr enum doesn't have " + arg2.Value())
	}

	return tokens.NewValueString("@parents "+arg0.Value()+"["+arg1.Value()+mod+"='"+arg2.Value()+"']", ctx), nil
}

func evalThisAttrEnumA(scope Scope, args []tokens.Token,
	ctx context.Context) (tokens.Token, error) {
	arg0, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	node := scope.GetNode()
	if _, ok := node.(*RootNode); ok {
		panic("cannot be rootnode")
	}

	if _, err := node.SearchAttrEnum(nil, arg0, ctx); err != nil {
		return nil, err
	}

	return tokens.NewValueString("@filter "+arg0.Value(), ctx), nil
}

func evalParentAttrEnumA(scope Scope, args []tokens.Token,
	ctx context.Context) (tokens.Token, error) {
	arg0, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	arg1, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	if _, err := scope.GetNode().SearchAttrEnum(arg0, arg1, ctx); err != nil {
		return nil, err
	}

	return tokens.NewValueString("@parents "+arg0.Value()+"["+arg1.Value()+"]", ctx), nil
}
