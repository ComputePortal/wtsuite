package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

// for debugging
const ELEMENT_COUNT = "__elementCount__"
const ELEMENT_COUNT_FOLDED = "__elementCountFolded__"

func evalElementCount(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {

	if len(args) != 0 && len(args) != 1 {
		return nil, ctx.NewError("Error: expected 0 or 1 arguments")
	}

	args, err := functions.EvalArgs(scope, args)
	if err != nil {
		return nil, err
	}

	folded := false
	if len(args) == 1 {
		b, err := tokens.AssertBool(args[0])
		if err != nil {
			return nil, err
		}

		folded = b.Value()
	}

	node := scope.GetNode()

	if folded {
		i := node.getElementCountFolded()
		return tokens.NewValueInt(i, ctx), nil
	} else {
		i := node.getElementCount()
		return tokens.NewValueInt(i, ctx), nil
	}
}
