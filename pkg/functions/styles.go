package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func minMaxScreenHeightWidth(args_ *tokens.Parens, isMin bool, isWidth bool, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	x, err := tokens.AssertFloat(args[0], "px")
	if err != nil {
		return nil, err
	}

	result := "@"
	if isMin {
		result += "min"
	} else {
		result += "max"
	}

	if isWidth {
		result += "-screen-width "
	} else {
		result += "-screen-height "
	}

	result += x.Write()

	return tokens.NewString(result, ctx)
}

func MaxScreenWidth(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, false, true, ctx)
}

func MaxScreenHeight(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, false, false, ctx)
}

func MinScreenWidth(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, true, true, ctx)
}

func MinScreenHeight(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, true, false, ctx)
}
