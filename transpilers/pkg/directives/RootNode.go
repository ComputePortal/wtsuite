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

func (n *RootNode) MapBlocks(blocks *tokens.StringDict) error {
	return nil
}

func (n *RootNode) GetOperations() []Operation {
	return []Operation{}
}

func (n *RootNode) PopOp(id string) (Operation, bool, error) {
	return nil, false, nil
}

func (n *RootNode) SearchStyle(key *tokens.String, ctx context.Context) (tokens.Token, error) {
	return tokens.NewNull(ctx), nil
}

func (n *RootNode) SearchAttrEnum(id *tokens.String, key *tokens.String, ctx context.Context) (*tokens.List, error) {
	if id == nil {
		return nil, ctx.NewError("Error: no parent element found")
	} else {
		return nil, ctx.NewError("Error: no parent element " + id.Value() + " found")
	}
}
