package directives

import (
	"io/ioutil"

	"../files"
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
)

func buildLinkedStyle(node Node, srcPath string, ctx context.Context) error {
	if style, err := tree.NewStyleSheetLink(srcPath, ctx); err != nil {
		return err
	} else {
		return node.AppendChild(style)
	}
}

func buildInlineStyle(node Node, attr *tokens.StringDict, content string,
	ctx context.Context) error {
	if style, err := tree.NewStyle(attr, content, ctx); err != nil {
		return err
	} else {
		return node.AppendChild(style)
	}
}

func Style(scope Scope, node Node, tag *tokens.Tag) error {
	ctx := tag.Context()

	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{"value"})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	srcToken_, hasSrc := attr.Get("src")
	if hasSrc && attr.Len() == 1 {
		if err := tag.AssertEmpty(); err != nil {
			context.AppendString(err, "Info: can't have both src and inline content")
			return err
		}

		srcToken, err := tokens.AssertString(srcToken_)
		if err != nil {
			return err
		}

		srcPath, err := files.Search(ctx.Path(), srcToken.Value())
		if err != nil {
			errCtx := srcToken.Context()
			return errCtx.NewError("Error: file " + err.Error())
		}

		if tree.INLINE {
			content, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return ctx.NewError(err.Error())
			}

			return buildInlineStyle(node, attr, string(content), ctx)
		} else {
			return buildLinkedStyle(node, srcPath, attr.Context())
		}
	} else {
		if !tag.IsEmpty() {
			errCtx := tag.Context()
			return errCtx.NewError("Error: inline style must use value attribute")
		}

		valueToken_, hasValue := attr.Get("value")
		if !hasValue {
			errCtx := attr.Context()
			return errCtx.NewError("Error: value attribute not found")
		}

		str, err := tokens.AssertString(valueToken_)
		if err != nil {
			return err
		}

		attr.Delete("value")
		return buildInlineStyle(node, attr, str.Value(), tag.Context())
	}
}

var _styleOk = registerDirective("style", Style)
