package directives

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type NodeType int

const (
	HTML NodeType = iota
	SVG
)

// styleTree and elementCount fit the node tree better than the scope tree
type Node interface {
  Parent() Node
	Name() string
	Type() NodeType

	getElementCount() int
	getElementCountFolded() int
	incrementElementCountFolded()
	getLastChild() tree.Tag
  getNode() Node

	AppendChild(tree.Tag) error

	PopOp(id string) (Operation, error) // for application
  GetBlockTarget(block *tokens.Tag) string // empty if not found
  SetBlockTarget(block *tokens.Tag, target string)

	SearchStyle(scope tokens.Scope, key *tokens.String, ctx context.Context) (tokens.Token, error)
  // register style sheet, which can be used by 
  RegisterStyleSheet(sheet StyleSheet) // the wraps of registered style sheets must be 
  // wraps after the main style sheet
}

type NodeData struct {
	tag    tree.Tag
	parent Node

	ecf int
}

func newNodeData(tag tree.Tag, parent Node) NodeData {
	return NodeData{tag, parent, 0}
}

func NewNode(tag tree.Tag, parent Node) *NodeData {
	node := newNodeData(tag, parent)
	return &node
}

func (n *NodeData) Parent() Node {
  return n.parent
}

func (n *NodeData) getNode() Node {
  return n
}

func (n *NodeData) GetBlockTarget(block *tokens.Tag) string {
  return n.parent.GetBlockTarget(block)
}

func (n *NodeData) SetBlockTarget(block *tokens.Tag, target string) {
  n.parent.SetBlockTarget(block, target)
}

func (n *NodeData) Name() string {
	return n.tag.Name()
}

func (n *NodeData) Type() NodeType {
	return n.parent.Type()
}

func (n *NodeData) Context() context.Context {
	return n.tag.Context()
}

func (n *NodeData) incrementElementCountFolded() {
	if n.Name() == "dummy" {
		n.parent.incrementElementCountFolded()
	} else {
		n.ecf += 1
	}
}

func (n *NodeData) getElementCountFolded() int {
	if n.tag != nil && n.Name() == "dummy" {
    if n.parent == nil {
      return 0
    } else {
      return n.parent.getElementCountFolded()
    }
	} else {
		return n.ecf
	}
}

func (n *NodeData) getElementCount() int {
  if (n.tag == nil) {
    return 0
  }

  return n.tag.NumChildren()
}

func (n *NodeData) getLastChild() tree.Tag {
	l := n.tag.NumChildren()
	if l == 0 {
		return nil
	} else {
		return n.tag.Children()[l-1]
	}
}

func (n *NodeData) AppendChild(child tree.Tag) error {
	n.tag.AppendChild(child)
	if child.Name() != "dummy" {
		n.incrementElementCountFolded()
	}

	return nil
}

func (n *NodeData) PopOp(id string) (Operation, error) {
	return n.parent.PopOp(id)
}

func (n *NodeData) SearchStyle(scope tokens.Scope, key *tokens.String, ctx context.Context) (tokens.Token, error) {
	attr := n.tag.Attributes()
	if styleToken_, ok := attr.Get("style"); ok && !tokens.IsNull(styleToken_) {
		styleToken, err := tokens.AssertStringDict(styleToken_)
		if err != nil {
			return nil, err
		}

		if v, ok := styleToken.Get(key.Value()); ok {
			return v, nil
		}
	}

  if n.parent == nil {
    if scope.Permissive() {
      return tokens.NewNull(ctx), nil
    } else {
      errCtx := key.Context()
      return nil, errCtx.NewError("Error: " + key.Value() + " not found in __pstyle__")
    }
  } else {
    return n.parent.SearchStyle(scope, key, ctx)
  }
}

func (n *NodeData) RegisterStyleSheet(sheet StyleSheet) {
  // parent can't be nil
  n.parent.RegisterStyleSheet(sheet)
}
