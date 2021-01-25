package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Fun interface {
	tokens.Token
	EvalFun(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error)
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
func NewFun(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	argsWithDefaults, err := tokens.AssertParens(args[0])
	if err != nil {
		return nil, err
	}

	// raw tokens should be ok, unless we want some special macro capabilities
  if err := argsWithDefaults.AssertUniqueNames(); err != nil {
    return nil, err
  }

	return NewAnonFun(scope, argsWithDefaults, args[1], ctx), nil
}

func EvalFun(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	args_, err := args_.EvalAsArgs(scope)
	if err != nil {
		return nil, err
	}

  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

  if len(args) != 2 {
    return nil, ctx.NewError("Error: expected 2 arguments")
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

  parens := tokens.NewParens(argList, nil, ctx)
  return fn.EvalFun(scope, parens, ctx)
}
