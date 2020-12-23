package functions

import (
	"math"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

// 3, or 5, or 7 etc. (last must always be else block
func IfElse(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if int(math.Mod(float64(len(args)-1), 2.0)) != 0 {
		return nil, ctx.NewError("Error: expected 3, 5, 7... arguments")
	}

	for i := 0; i < len(args); i += 2 {
		if i < len(args)-1 {
			argCond, err := args[i].Eval(scope)
			if err != nil {
				return nil, err
			}

			cond, err := tokens.AssertBool(argCond)
			if err != nil {
				return nil, err
			}

			if cond.Value() {
				return args[i+1].Eval(scope)
			}
		}
	}

	return args[len(args)-1].Eval(scope)
}
