package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func getUndefinedOrNull(scope Scope, keyToken_ tokens.Token, ctx context.Context) (tokens.Token, error) {
	keyToken, err := keyToken_.Eval(scope)
	if err != nil {
		return nil, err
	}

	nameToken, err := tokens.AssertString(keyToken)
	if err != nil {
		return nil, err
	}

	key := nameToken.Value()

	var res tokens.Token = nil
	switch {
	case HasDefine(key):
		res = GetDefine(key)
	case scope.HasVar(key):
		res = scope.GetVar(key).Value
	case key == URL:
		res, _ = GetActiveURL(ctx)
	case functions.HasFun(key):
		res = functions.NewBuiltInFun(key, ctx)
	}

	return res, nil
}

// like get but also returns backup if null
func evalVar(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	arg0 := args[0]
	if tokens.IsFunction(args[0], "get") {
		fn, err := tokens.AssertFunction(args[0])
		if err != nil {
			panic(err)
		}

		fnArgs := fn.Args()
		if len(fnArgs) == 1 {
			arg0, err = getUndefinedOrNull(scope, fnArgs[0], ctx)
			if err != nil {
				return nil, err
			}
		} else {
			arg0, err = arg0.Eval(scope)
			if err != nil {
				return nil, err
			}
		}
	} else {
		var err error
		arg0, err = arg0.Eval(scope)
		if err != nil {
			return nil, err
		}
	}

	if arg0 == nil || tokens.IsNull(arg0) {
		return args[1].Eval(scope)
	} else {
		return arg0, nil
	}
}
