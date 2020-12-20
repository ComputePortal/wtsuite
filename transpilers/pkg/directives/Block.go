package directives

import (
  tokens "../tokens/html"
)

func EvalBlock(scope Scope, node Node, tag *tokens.Tag) error {
  opName, err := getOpNameTarget("name", tag)
  if err != nil {
    return err
  }

  if !InsideTemplateScope(scope) {
    errCtx := tag.Context()
    return errCtx.NewError("Error: block not inside template")
  }

  if !IsUniqueOpTargetName(opName) {
    panic("should've been deferred and given a unique name earlier")
  }

  // first create self as a dummy node
  dummyTag := tokens.NewTag("dummy", tokens.NewEmptyRawDict(tag.Context()), tag.Children(),
    tag.Context())

  if err := buildTree(scope, node, node.Type(), dummyTag, opName); err != nil {
    return err
  }

  return nil
}

var _evalBlockOk = registerDirective("block", EvalBlock)
