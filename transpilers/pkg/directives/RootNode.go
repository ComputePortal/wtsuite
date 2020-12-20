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

func (n *RootNode) GetOperations() []Operation {
	return []Operation{}
}

func (n *RootNode) PopOp(id string) (Operation, error) {
	return nil, nil
}

func (n *RootNode) SearchStyle(key *tokens.String, ctx context.Context) (tokens.Token, error) {
	return tokens.NewNull(ctx), nil
}
