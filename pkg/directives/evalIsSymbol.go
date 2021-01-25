package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func evalIsSymbol(scope Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  var err error
  args_, err = args_.EvalAsArgs(scope)
  if err != nil {
    return nil, err
  }

  args, err := functions.CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

  arg, err := tokens.AssertString(args[0])
  if err != nil {
    return nil, err
  }

  b := scope.HasVar(arg.Value()) || scope.HasTemplate(arg.Value())

  return tokens.NewValueBool(b, ctx), nil
}
