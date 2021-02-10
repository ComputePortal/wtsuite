package directives

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type RootNode struct {
	t NodeType
  sheets []StyleSheet
	NodeData
}

func NewRootNode(tag tree.Tag, t NodeType) *RootNode {
	if tag.Name() != "" {
		panic("expected Root or SVGRoot (tags with empty names)")
	}

	return &RootNode{t, make([]StyleSheet, 0), newNodeData(tag, nil)}
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

func (n *RootNode) RegisterStyleSheet(sheet StyleSheet) {
  n.sheets = append(n.sheets, sheet)
}
