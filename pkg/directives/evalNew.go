package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func evalNew(scope Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  var err error
	args_, err = args_.EvalAsArgs(scope)
	if err != nil {
		return nil, err
	}

  args, err := functions.CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}


	nameToken, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	if err := AssertValidVar(nameToken); err != nil {
		return nil, err
	}

	valueToken := args[1]

	key := nameToken.Value()
	switch {
	case HasGlobal(key):
		errCtx := nameToken.InnerContext()
		return nil, errCtx.NewError("Error: can't redefine global")
	case scope.HasVar(key):
		errCtx := nameToken.InnerContext()
		return nil, errCtx.NewError("Error: can't redefine variable")
	default:
		v := functions.Var{valueToken, false, false, false, false, ctx}
    if err := scope.SetVar(key, v); err != nil {
      return nil, err
    }
	}

	return valueToken, nil
}
