package directives

import (
	"github.com/computeportal/wtsuite/pkg/tree"
)

type SVGNode struct {
	NodeData
}

func NewSVGNode(tag tree.Tag, parent Node) *SVGNode {
	return &SVGNode{newNodeData(tag, parent)}
}

func (n *SVGNode) Type() NodeType {
	return SVG
}

func (n *SVGNode) getNode() Node {
  return n
}
