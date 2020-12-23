package directives

import (
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

// scripts are usefull for internet explorer fallback
func Script(scope Scope, node Node, tag *tokens.Tag) error {
	if !tag.IsScript() || tag.Name() != "script" {
		panic("expected js script tag")
	}

	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{"value"})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	if idToken, ok := attr.Get("id"); ok {
		errCtx := idToken.Context()
		return errCtx.NewError("Error: id for scripts is meaningless")
	}

	content := ""
	_, ok := attr.Get("value")
	if ok {
		valueToken, err := tokens.DictString(attr, "value")
		if err != nil {
			return err
		}

		content = valueToken.Value()
	}

	attr.Delete("value")

	script, err := tree.NewScript(attr, content, tag.Context())
	if err != nil {
		return err
	}

	return node.AppendChild(script)
}

var _scriptOk = registerDirective("script", Script)
