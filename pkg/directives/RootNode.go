package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
)

type RootNode struct {
	t NodeType
	NodeData
}

func NewRootNode(tag tree.Tag, t NodeType) *RootNode {
	if tag.Name() != "" {
		panic("expected Root or SVGRoot (tags with empty names)")
	}

	return &RootNode{t, newNodeData(tag, nil)}
}

func (n *RootNode) Type() NodeType {
	return n.t
}

func (n *RootNode) PopOp(id string) (Operation, error) {
	return nil, nil
}

func (n *RootNode) SearchStyle(scope tokens.Scope, key *tokens.String, ctx context.Context) (tokens.Token, error) {
  if scope.Permissive() {
    return tokens.NewNull(ctx), nil
  } else {
    return nil, ctx.NewError("Error: key " + key.Value() + " not found in __pstyle__")
  }
}

func (n *RootNode) SetBlockTarget(block *tokens.Tag, target string) {
}

func (n *RootNode) GetBlockTarget(block *tokens.Tag) string {
  return ""
}

func (n *RootNode) getNode() Node {
  return n
}
