package directives

import (
	tokens "../tokens/html"
	"../tree"
)

type URINode struct {
  blockTargets map[*tokens.Tag]string
	NodeData
}

// a URINode doesnt have a parent
func NewURINode() *URINode {
	return &URINode{make(map[*tokens.Tag]string), newNodeData(nil, nil)}
}

func (n *URINode) GetBlockTarget(block *tokens.Tag) string {
  opName, ok := n.blockTargets[block]
  if !ok {
    return ""
  }

  return opName
}

func (n *URINode) SetBlockTarget(block *tokens.Tag, target string) {
  n.blockTargets[block] = target
}

func (n *URINode) PopOp(id string) (Operation, error) {
	return nil, nil
}

func (n *URINode) incrementElementCountFolded() {
	n.ecf += 1
}

func (n *URINode) Type() NodeType {
  return SVG
}

func (n *URINode) getNode() Node {
  return n
}

func (n *URINode) AppendChild(tag tree.Tag) error {
	if n.tag != nil {
		errCtx := tag.Context()
		return errCtx.NewError("Error: unexpected second tag")
	}

	n.tag = tag

	return nil
}
