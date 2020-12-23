package functions

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Replace(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 3 {
		return nil, ctx.NewError("Error: expected 3 arguments")
	}

	str, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	old, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	new, err := tokens.AssertString(args[2])
	if err != nil {
		return nil, err
	}

	replaced := strings.Replace(str.Value(), old.Value(), new.Value(), -1)

	return tokens.NewValueString(replaced, ctx), nil
}
