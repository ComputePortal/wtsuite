package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func evalGet(scope Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  var err error
  args_, err = args_.EvalAsArgs(scope)
	if err != nil {
		return nil, err
	}

  args, err := functions.CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	var fallback tokens.Token = nil
	if len(args) == 2 {
		fallback = args[1]
	} else if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}

	nameToken, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	key := nameToken.Value()
	switch {
	case HasDefine(key):
		return GetDefine(key), nil
	case scope.HasVar(key): // prefer variable over builtin function
		return scope.GetVar(key).Value, nil
	case key == URL:
		return GetActiveURL(ctx)
	case functions.HasFun(key):
		return functions.NewBuiltInFun(key, ctx), nil
	case fallback != nil:
		return fallback, nil
	default:
		errCtx := nameToken.InnerContext()
		err := errCtx.NewError("Error: variable '" + key + "' not defined")
		if key == ELEMENT_COUNT {
			context.AppendString(err, "Hint: "+ELEMENT_COUNT+" is only available inside tags")
		}

		context.AppendString(err, "Info: available names\n"+scope.listValidVarNames())
		return nil, err
	}
}
