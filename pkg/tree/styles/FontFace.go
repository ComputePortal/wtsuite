package styles

import (
	"../../tokens/context"
	tokens "../../tokens/html"
)

func ImportFontFace(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments for @font-face")
	}

	if tokens.IsNull(v) {
		return []Rule{}, nil
	}

	attr, err := tokens.AssertStringDict(v)
	if err != nil {
		return nil, err
	}

	leafAttr, rules, err := expandNested(attr, sel)
	if err != nil {
		return nil, err
	}

	if len(rules) != 0 {
		return nil, ctx.NewError("Error: unexpected nested rules")
	}

	return []Rule{NewSelectorRule(NewSelector("@font-face"), leafAttr)}, nil
}

var _importFontFaceOk = registerAtFunction("font-face", ImportFontFace)
