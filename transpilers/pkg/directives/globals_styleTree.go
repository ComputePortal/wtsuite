package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

func evalSearchStyle(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 && len(args) != 2 {
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}

	args, err := functions.EvalArgs(scope, args)
	if err != nil {
		return nil, err
	}

	key, err := tokens.AssertString(args[len(args)-1])
	if err != nil {
		return nil, err
	}

	if len(args) == 2 && !tokens.IsNull(args[0]) {
		d, err := tokens.AssertStringDict(args[0])
		if err != nil {
			return nil, err
		}

		if val, ok := d.Get(key.Value()); ok {
			return val, nil
		}
	}

	node := scope.GetNode()

	return node.SearchStyle(key, ctx)
}

func SearchStyle(scope Scope, tagAttr *tokens.StringDict, key string, ctx context.Context) (tokens.Token, error) {
	if styleToken_, ok := tagAttr.Get("style"); ok && !tokens.IsNull(styleToken_) {
		styleToken, err := tokens.AssertStringDict(styleToken_)
		if err != nil {
			return nil, err
		}

		if v, ok := styleToken.Get(key); ok {
			return v, nil
		}
	}

	val, err := evalSearchStyle(scope, []tokens.Token{tokens.NewValueString(key, ctx)}, ctx)

	return val, err
}
