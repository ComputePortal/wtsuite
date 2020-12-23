package directives

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type TemplateNode struct {
	parent     Node
	operations []Operation
  collectDeferred bool
  blockTargets map[*tokens.Tag]string
  ctx context.Context
}

func NewTemplateNode(parent Node, ctx context.Context) *TemplateNode {
	return &TemplateNode{parent, make([]Operation, 0), false, make(map[*tokens.Tag]string), ctx}
}

func (n *TemplateNode) Parent() Node {
  return n.parent
}

func (n *TemplateNode) Name() string {
	return n.parent.Name()
}

func (n *TemplateNode) Type() NodeType {
	return n.parent.Type()
}

func (n *TemplateNode) getNode() Node {
  return n.parent
}

func (n *TemplateNode) getElementCount() int {
	return n.parent.getElementCount()
}

func (n *TemplateNode) getElementCountFolded() int {
	return n.parent.getElementCountFolded()
}

func (n *TemplateNode) incrementElementCountFolded() {
	n.parent.incrementElementCountFolded()
}

func (n *TemplateNode) getLastChild() tree.Tag {
	return n.parent.getLastChild()
}

func (n *TemplateNode) AppendChild(child tree.Tag) error {
	return n.parent.AppendChild(child)
}

func (n *TemplateNode) SearchStyle(scope tokens.Scope, key *tokens.String, ctx context.Context) (tokens.Token, error) {
	return n.parent.SearchStyle(scope, key, ctx)
}

func (n *TemplateNode) StartDeferral() {
  n.collectDeferred = true
}

func (n *TemplateNode) StopDeferral() {
  n.collectDeferred = false
}

func IsDeferringTemplateNode(node Node) bool {
  if tNode, ok := node.(*TemplateNode); ok && tNode.collectDeferred {
    return true
  } else {
    return false
  }
}

func (n *TemplateNode) AssertAllOperationsDone(ctx context.Context) error {
  var err error = nil
  for _, op := range n.operations {
    if !op.Popped() {
      if err == nil {
        err = ctx.NewError("Error: unapplied ops")
      }

      context.AppendContextString(err, "Info: not applied to "+op.Target(), op.Context())
    }
  }

  return err
}

func (n *TemplateNode) PopOp(target string) (Operation, error) {
	parentOp, err := n.parent.PopOp(target)
	if err != nil {
		return nil, err
	}

	var thisOp Operation = nil
	thisOk := false
	for i, op := range n.operations {
		if op.Target() == target {
      if !IsUniqueOpTargetName(target) {
        // shared ops are not removed
        if i < len(n.operations)-1 {
          n.operations = append(n.operations[0:i], n.operations[i+1:]...)
        } else {
          n.operations = n.operations[0:i]
        }
      }

      op.SetPopped()

			thisOp = op
			thisOk = true
			break
		}
	}

	if thisOk && parentOp != nil {
		// merge
		mop, err := thisOp.Merge(parentOp)
		if err != nil {
			return nil, err
		}

    mop.SetPopped()

		return mop, nil
	} else if thisOk {
		return thisOp, nil
	} else {
		return parentOp, nil
	}
}

func (n *TemplateNode) PushOp(op Operation) error {
  // no other ops can have this name
  for i, prev := range n.operations {
    if prev.Target() == op.Target() {
      var err error
      n.operations[i], err = prev.Merge(op)
      if err != nil {
        return err
      }
      return nil
    }
  }

  n.operations = append(n.operations, op)

  return nil
}

func (n *TemplateNode) AppendToDefault(scope Scope, tag *tokens.Tag) error {
  appToDef, err := NewAppendToDefaultOp(scope, []*tokens.Tag{tag})
  if err != nil {
    return err
  }

  // merge with any previous 
  for i, prev := range n.operations {
    if prev.Target() == "default" {
      var err error
      n.operations[i], err = prev.Merge(appToDef)
      if err != nil {
        return err
      }

      return nil
    }
  }

  n.operations = append(n.operations, appToDef)
  return nil
}

func (n *TemplateNode) SetBlockTarget(block *tokens.Tag, target string) {
  n.blockTargets[block] = target
}

func (n *TemplateNode) GetBlockTarget(block *tokens.Tag) string {
  opName, ok := n.blockTargets[block]
  if !ok {
    return n.parent.GetBlockTarget(block)
  }

  return opName
}

func GetTemplateNode(node Node) *TemplateNode {
  if tNode, ok := node.(*TemplateNode); ok {
    return tNode
  }

  if parent := node.Parent(); parent != nil {
    return GetTemplateNode(parent)
  }

  return nil
}
