package directives

import (
	"../functions"
	tokens "../tokens/html"
)

func For(scope Scope, node Node, tag *tokens.Tag) error {
	ctx := tag.Context()

	subScope := NewStatementScope(scope)

	attr, err := tag.Attributes([]string{"iname", "vname"})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	valuesToken, err := tokens.DictList(attr, "in")
	if err != nil {
		return err
	}
	argCount := 1

	var inameToken *tokens.String = nil
	var vnameToken *tokens.String = nil

	vnameToken_, hasVName := attr.Get("vname")
	if inameToken_, ok := attr.Get("iname"); ok {
		argCount++
		if !hasVName {
			vnameToken, err = tokens.AssertString(inameToken_)
			if err != nil {
				return err
			}
		} else {
			inameToken, err = tokens.AssertString(inameToken_)
			if err != nil {
				return err
			}
		}
	}

	if hasVName {
		argCount++
		vnameToken, err = tokens.AssertString(vnameToken_)
		if err != nil {
			return err
		}
	}

	if attr.Len() != argCount {
		errCtx := attr.Context()
		return errCtx.NewError("Error: unexpected attributes")
	}

	if err := valuesToken.Loop(func(i int, v tokens.Token, last bool) error {
		if inameToken != nil {
			iVar := functions.Var{tokens.NewValueInt(i, ctx), true, true, false, false, ctx}
			subScope.SetVar(inameToken.Value(), iVar)
		}
		if vnameToken != nil {
			vVar := functions.Var{v, true, true, false, false, ctx}
			subScope.SetVar(vnameToken.Value(), vVar)
		}

		for _, child := range tag.Children() {
			if err := BuildTag(subScope, node, child); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

var _forOk = registerDirective("for", For)
