package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func evalGet(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	// avoid false wasWord positives
	if str, ok := args[0].(*tokens.String); ok && str.WasWord() {
		args[0] = tokens.NewValueString(str.Value(), str.Context())
	}

	args, err := functions.EvalArgs(scope, args)
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
