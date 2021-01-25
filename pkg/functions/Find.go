package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func findWithFunction(scope tokens.Scope, arg0 *tokens.List, fn Fun, ctx context.Context) (tokens.Token, error) {
  innerArgs := make([]tokens.Token, fn.Len())
  result := -1

  if err := arg0.Loop(func(i int, value tokens.Token, last bool) error {
    if result != -1 {
      return nil
    }

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

    cond, err := fn.EvalFun(scope, tokens.NewParens(innerArgs, nil, ctx), ctx)
    if err != nil {
      return err
    }

    b, err := tokens.AssertBool(cond)
    if err != nil {
      return err
    }

    if b.Value() {
      result = i
    }

    return nil
  }); err != nil {
    return nil, err
  }

  if result == -1 {
    return tokens.NewNull(ctx), nil
  } else {
    return tokens.NewValueInt(result, ctx), nil
  }
}

func findWithIsSame(scope tokens.Scope, arg0 *tokens.List, arg1 tokens.Token, ctx context.Context) (tokens.Token, error) {
  iRes := -1

  if err := arg0.Loop(func(i int, value tokens.Token, last bool) error {
    if iRes != -1 {
      return nil
    }

    if arg1.IsSame(value) {
      iRes = i
    }

    return nil
  }); err != nil {
    return nil, err
  }

  if iRes == -1 {
    return tokens.NewNull(ctx), nil
  } else {
    return tokens.NewValueInt(iRes, ctx), nil
  }
}

func Find(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

  if len(args) != 2 {
    return nil, ctx.NewError("Error: expected 2 arguments")
  }

  arg0_, err := args[0].Eval(scope)
  if err != nil {
    return nil, err
  }

  arg0, err := tokens.AssertList(arg0_)
  if err != nil {
    return nil, err
  }

	arg1, err := args[1].Eval(scope)
	if err != nil {
		return nil, err
	}

  if IsFun(arg1) {
    fn, err := AssertFun(arg1)
    if err != nil {
      panic(err)
    }

    return findWithFunction(scope, arg0, fn, ctx)
  } else {
    arg1, err := args[1].Eval(scope)
    if err != nil {
      return nil, err
    }

    return findWithIsSame(scope, arg0, arg1, ctx)
  }
}
