package directives

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type BuildStyleFunc func(d *tokens.StringDict) (string, error)
var BuildStyle BuildStyleFunc = nil
func RegisterBuildStyle(fn BuildStyleFunc) bool {
  BuildStyle = fn
  return true
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
	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	contentToken, ok := attr.Get(".content")
  if !ok {
    errCtx := tag.Context()
    return errCtx.NewError("Error: style without content")
  }

  // build the style
  d, err := tokens.AssertStringDict(contentToken)
  if err != nil {
    return err
  }

  contentStr, err := BuildStyle(d)
  if err != nil {
    return err
  }

  attr.Delete(".content")

  return buildInlineStyle(node, attr, contentStr, tag.Context())
}

var _styleOk = registerDirective("style", Style)
