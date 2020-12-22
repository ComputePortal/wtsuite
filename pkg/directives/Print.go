package directives

import (
	tokens "../tokens/html"
	"../tree"
)

func Print(scope Scope, node Node, tag *tokens.Tag) error {
	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{"value"})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	valueToken, err := tokens.DictPrimitive(attr, "value")
	if err != nil {
		return err
	}

	if err := tag.AssertEmpty(); err != nil {
		return err
	}

	return node.AppendChild(tree.NewText(valueToken.Write(), valueToken.Context()))
}

var _printOk = registerDirective("print", Print)
