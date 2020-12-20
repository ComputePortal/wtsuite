package directives

import (
	"../tree"
)

type URINode struct {
	NodeData
}

func NewURINode(parent Node) *URINode {
	return &URINode{newNodeData(nil, parent)}
}

func (n *URINode) PopOp(id string) (Operation, error) {
	return nil, nil
}

func (n *URINode) incrementElementCountFolded() {
	n.ecf += 1
}

func (n *URINode) AppendChild(tag tree.Tag) error {
	if n.tag != nil {
		errCtx := tag.Context()
		return errCtx.NewError("Error: unexpected second tag")
	}

	n.tag = tag

	return nil
}
