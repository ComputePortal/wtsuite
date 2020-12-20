package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
)

type TemplateNode struct {
	parent     Node
	operations []Operation
  collectDeferred bool
}

func NewTemplateNode(parent Node) *TemplateNode {
	return &TemplateNode{parent, make([]Operation, 0), false}
}

func (n *TemplateNode) Name() string {
	return n.parent.Name()
}

func (n *TemplateNode) Type() NodeType {
	return n.parent.Type()
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

func (n *TemplateNode) SearchStyle(key *tokens.String, ctx context.Context) (tokens.Token, error) {
	return n.parent.SearchStyle(key, ctx)
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

func (n *TemplateNode) GetOperations() []Operation {
	return n.operations
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
			if i < len(n.operations)-1 {
				n.operations = append(n.operations[0:i], n.operations[i+1:]...)
			} else {
				n.operations = n.operations[0:i]
			}

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
