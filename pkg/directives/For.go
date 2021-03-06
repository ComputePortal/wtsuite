package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func For(scope Scope, node Node, tag *tokens.Tag) error {
	ctx := tag.Context()

	subScope := NewBranchScope(scope)

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
    loopScope := NewBranchScope(subScope)

		if inameToken != nil {
			iVar := functions.Var{tokens.NewValueInt(i, ctx), true, true, false, false, ctx}
      if err := loopScope.SetVar(inameToken.Value(), iVar); err != nil {
        return err
      }
		}
		if vnameToken != nil {
			vVar := functions.Var{v, true, true, false, false, ctx}
      if err := loopScope.SetVar(vnameToken.Value(), vVar); err != nil {
        return err
      }
		}

		for _, child := range tag.Children() {
			if err := BuildTag(loopScope, node, child); err != nil {
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
