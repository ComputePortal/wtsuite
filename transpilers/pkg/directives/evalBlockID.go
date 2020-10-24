package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func evalBlockID(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected exactly one argument")
	}

	arg, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	if !scope.HasVar("__blocks__") {
		return nil, ctx.NewError("Error: no block-id defined in this scope")
	} else {
		blocksVar := scope.GetVar("__blocks__")

		blocksVal, err := tokens.AssertStringDict(blocksVar.Value)
		if err != nil {
			panic(err)
		}

		idToken_, hasArg := blocksVal.Get(arg.Value())
		if !hasArg {
			if scope.Parent() == nil {
				return nil, ctx.NewError("Error: block id '" + arg.Value() + "' not found")
			} else {
				return evalBlockID(scope.Parent(), args, ctx)
			}
		} else {
			idToken, err := tokens.AssertString(idToken_)
			if err != nil {
				panic(err)
			}

			return idToken, nil
		}
	}
}
