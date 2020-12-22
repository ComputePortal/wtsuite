package styles

import (
	"strings"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type MinMaxScreenRangeType int

const (
	MIN_SCREEN MinMaxScreenRangeType = 1 << 1
	MAX_SCREEN                       = 1 << 2
)

type Media struct {
	Query
}

func minMaxScreenHeightWidth(sel Selector, args []string, v tokens.Token,
	r MinMaxScreenRangeType, isWidth bool, ctx context.Context) ([]Rule, error) {
	if (r&MIN_SCREEN) > 0 && (r&MAX_SCREEN) > 0 {
		if len(args) != 2 {
			return nil, ctx.NewError("Error: expected 2 arguments for screen")
		}
	} else if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	widthOrHeight := "-height"
	if isWidth {
		widthOrHeight = "-width"
	}

	if tokens.IsNull(v) {
		return []Rule{}, nil
	}

	var b strings.Builder
	b.WriteString("@media only screen and")

	argMaxI := 0
	if (r & MIN_SCREEN) > 0 {
		b.WriteString(" (min")
		b.WriteString(widthOrHeight)
		b.WriteString(":")
		b.WriteString(args[0])
		b.WriteString(")")

		if (r & MAX_SCREEN) > 0 {
			b.WriteString(" and")
			argMaxI = 1
		}
	}

	if (r & MAX_SCREEN) > 0 {
		b.WriteString(" (max")
		b.WriteString(widthOrHeight)
		b.WriteString(":")
		b.WriteString(args[argMaxI])
		b.WriteString(")")
	}

	q, err := newQuery(sel, v, b.String())
	if err != nil {
		return nil, err
	}

	return []Rule{&Media{q}}, nil
}

func MaxScreenWidth(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return minMaxScreenHeightWidth(sel, args, v, MAX_SCREEN, true, ctx)
}

func MaxScreenHeight(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return minMaxScreenHeightWidth(sel, args, v, MAX_SCREEN, false, ctx)
}

func MinScreenWidth(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return minMaxScreenHeightWidth(sel, args, v, MIN_SCREEN, true, ctx)
}

func MinScreenHeight(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return minMaxScreenHeightWidth(sel, args, v, MIN_SCREEN, false, ctx)
}

func ScreenWidth(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return minMaxScreenHeightWidth(sel, args, v, MIN_SCREEN|MAX_SCREEN, true, ctx)
}

func ScreenHeight(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return minMaxScreenHeightWidth(sel, args, v, MIN_SCREEN|MAX_SCREEN, false, ctx)
}

var _maxScreenHeightOk = registerAtFunction("max-screen-height", MaxScreenHeight)
var _minScreenHeightOk = registerAtFunction("min-screen-height", MinScreenHeight)
var _maxScreenWidthOk = registerAtFunction("max-screen-width", MaxScreenWidth)
var _minScreenWidthOk = registerAtFunction("min-screen-width", MinScreenWidth)
var _screenHeightOk = registerAtFunction("screen-height", ScreenHeight)
var _screenWidthOk = registerAtFunction("screen-width", ScreenWidth)
