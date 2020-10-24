package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func minMaxScreenHeightWidth(args []tokens.Token, isMin bool, isWidth bool, ctx context.Context) (tokens.Token, error) {
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

func MaxScreenWidth(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, false, true, ctx)
}

func MaxScreenHeight(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, false, false, ctx)
}

func MinScreenWidth(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, true, true, ctx)
}

func MinScreenHeight(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return minMaxScreenHeightWidth(args, true, false, ctx)
}
