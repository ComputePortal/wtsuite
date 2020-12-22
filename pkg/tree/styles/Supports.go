package styles

import (
	"strings"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type Supports struct {
	Query
}

func browserFamilyOnly(sel Selector, args []string, featureStr string, v tokens.Token,
	ctx context.Context) ([]Rule, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	var b strings.Builder

	b.WriteString("@supports ")
	b.WriteString(featureStr)

	q, err := newQuery(sel, v, b.String())
	if err != nil {
		return nil, err
	}

	return []Rule{&Supports{q}}, nil
}

func MozOnly(sel Selector, args []string, v tokens.Token,
	ctx context.Context) ([]Rule, error) {
	return browserFamilyOnly(sel, args, "(-moz-appearance:meterbar)", v, ctx)
}

func WebkitOnly(sel Selector, args []string, v tokens.Token,
	ctx context.Context) ([]Rule, error) {
	return browserFamilyOnly(sel, args, "(-webkit-appearance:none)", v, ctx)
}

func MSIEOnly(sel Selector, args []string, v tokens.Token,
	ctx context.Context) ([]Rule, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	var b strings.Builder

	b.WriteString("@media screen and (min-width:0\\0)")

	q, err := newQuery(sel, v, b.String())
	if err != nil {
		return nil, err
	}

	return []Rule{&Media{q}}, nil
}

var _mozOnlyOk = registerAtFunction("moz-only", MozOnly)
var _webkitOnlyOk = registerAtFunction("webkit-only", WebkitOnly)
var _msieOnlyOk = registerAtFunction("msie-only", MSIEOnly)
