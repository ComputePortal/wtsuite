package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Fun interface {
	tokens.Token
	EvalFun(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error)
	Len() int // number of arguments, -1: variable
}

func IsFun(t tokens.Token) bool {
	_, ok := t.(Fun)
	return ok
}

func AssertFun(t tokens.Token) (Fun, error) {
	f, ok := t.(Fun)
	if ok {
		return f, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected function")
	}
}

// eg. function([args..], [body1;body2;body3][-1])
func NewFun(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	argsWithDefaults, err := tokens.AssertParens(args[0])
	if err != nil {
		return nil, err
	}

	// raw tokens should be ok, unless we want some special macro capabilities

	//arg0, err := args[0].Eval(scope)
	//if err != nil {
	//return nil, err
	//}

	// list can only contain names
	funArgNames := make([]string, argsWithDefaults.Len())
	checkUniqueness := make(map[string]*tokens.String)
	if err := argsWithDefaults.Loop(func(i int, argName tokens.Token, argDef tokens.Token) error {
		argNameToken, err := tokens.AssertString(argName)
		if err != nil {
			return err
		}

		if prev, ok := checkUniqueness[argNameToken.Value()]; ok {
			errCtx := context.MergeContexts(prev.Context(), argNameToken.Context())
			return errCtx.NewError("Error: duplicate arg names")
		}

		checkUniqueness[argNameToken.Value()] = argNameToken

		funArgNames[i] = argNameToken.Value()
		return nil
	}); err != nil {
		return nil, err
	}

	return NewAnonFun(scope, funArgNames, argsWithDefaults.Alts(), args[1], ctx), nil
}

func EvalFun(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	args, err := EvalArgs(scope, args)
	if err != nil {
		return nil, err
	}

	fn, err := AssertFun(args[0])
	if err != nil {
		return nil, err
	}

	// a list is better than varargs, because it can be processed by builtin list functions
	list, err := tokens.AssertList(args[1])
	if err != nil {
		return nil, err
	}

	argList := make([]tokens.Token, list.Len())
	if err := list.Loop(func(i int, v tokens.Token, last bool) error {
		argList[i] = v
		return nil
	}); err != nil {
		panic(err)
	}

	return fn.EvalFun(scope, argList, ctx)
}
