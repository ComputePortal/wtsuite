package directives
 
import (

  tokens "../tokens/html"
)

func EvalBlock(scope Scope, node Node, tag *tokens.Tag) error {
  tNode := GetTemplateNode(node) 
  if tNode == nil {
    errCtx := tag.Context()
    return errCtx.NewError("Error: block not inside template")
  }

  opName := tNode.GetBlockTarget(tag)

  if opName == "" {
    panic("not found in block targets")
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
