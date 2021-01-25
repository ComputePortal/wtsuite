package directives

import (
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

func TOC(scope Scope, node Node, tag *tokens.Tag) error {
  ctx := tag.Context()

  attrScope := NewSubScope(scope)

  setLazyTagVars(attrScope, ctx)
 
  if err := tag.AssertEmpty(); err != nil {
    return err
  }

  attr, err := tag.Attributes([]string{"numbering"})
  if err != nil {
    return err
  }

	attr, err = attr.EvalStringDict(attrScope)
	if err != nil {
		return err
	}

  attr, err = removeForcedSuffix(attr)
  if err != nil {
    return err
  }

  numbering := "none"
  numberingToken_, ok := attr.Get("numbering")
  if ok {
    numberingToken, err := tokens.AssertString(numberingToken_)
    if err != nil {
      return err
    }

    numbering = numberingToken.Value()
    attr.Delete("numbering")

    if numbering != "roman" && numbering != "none" && numbering != "decimal" {
      errCtx := numberingToken_.Context()
      return errCtx.NewError("Error: invalid numbering scheme (expected \"decimal\", \"roman\" or \"none\"")
    }
  }

  tocTag, err := tree.NewTOC(numbering, attr, ctx)
  if err != nil {
    return err
  }

  if err := node.AppendChild(tocTag); err != nil {
    return nil
  }

  return nil
}

var tocOK = registerDirective("toc", TOC)
