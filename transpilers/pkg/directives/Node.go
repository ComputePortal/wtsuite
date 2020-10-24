package directives

import (
	"../tokens/context"
	tokens "../tokens/html"

	"../tree"
)

type NodeType int

const (
	HTML NodeType = iota
	SVG
)

// styleTree and elementCount fit the node tree better than the scope tree
type Node interface {
	Name() string
	Type() NodeType

	getElementCount() int
	getElementCountFolded() int
	incrementElementCountFolded()
	getLastChild() tree.Tag

	AppendChild(tree.Tag) error

	MapBlocks(blocks *tokens.StringDict) error
	GetOperations() []Operation               // for merging
	PopOp(id string) (Operation, bool, error) // for application

	SearchStyle(key *tokens.String, ctx context.Context) (tokens.Token, error)
	// if id == nil search this node only
	SearchAttrEnum(id *tokens.String, key *tokens.String, ctx context.Context) (*tokens.List, error)
}

type NodeData struct {
	tag    tree.Tag
	parent Node

	ecf int
	// TODO: list of operations
}

func newNodeData(tag tree.Tag, parent Node) NodeData {
	return NodeData{tag, parent, 0}
}

func NewNode(tag tree.Tag, parent Node) *NodeData {
	node := newNodeData(tag, parent)
	return &node
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
	if n.Name() == "dummy" {
		return n.parent.getElementCountFolded()
	} else {
		return n.ecf
	}
}

func (n *NodeData) getElementCount() int {
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
	attr := child.Attributes()
	if attr != nil {
		// for debugging
		attr.Set(ELEMENT_COUNT, tokens.NewValueInt(n.getElementCount(), child.Context()))
		attr.Set(ELEMENT_COUNT_FOLDED, tokens.NewValueInt(n.getElementCountFolded(),
			child.Context()))
	}

	n.tag.AppendChild(child)
	if child.Name() != "dummy" {
		n.incrementElementCountFolded()
	}
	return nil
}

func (n *NodeData) MapBlocks(blocks *tokens.StringDict) error {
	return n.parent.MapBlocks(blocks)
}

func (n *NodeData) GetOperations() []Operation {
	return n.parent.GetOperations()
}

func (n *NodeData) PopOp(id string) (Operation, bool, error) {
	return n.parent.PopOp(id)
}

func (n *NodeData) SearchStyle(key *tokens.String, ctx context.Context) (tokens.Token, error) {
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

	return n.parent.SearchStyle(key, ctx)
}

func (n *NodeData) SearchAttrEnum(id *tokens.String, key *tokens.String, ctx context.Context) (*tokens.List, error) {
	if id != nil && (id.Value() != n.tag.GetID() && !(id.Value() == "html" && n.tag.Name() == "html") && !(id.Value() == "body" && n.tag.Name() == "body")) {
		return n.parent.SearchAttrEnum(id, key, ctx)
	}

	attr := n.tag.Attributes()

	if v_, ok := attr.Get(key.Value()); !ok {
		return nil, ctx.NewError("Error: attr " + key.Value() + " not found")
	} else {
		if !tokens.IsAttrEnumList(v_) {
			return nil, ctx.NewError("Error: attr " + key.Value() + " value is not an enum list")
		}

		lst, err := tokens.AssertList(v_)
		if err != nil {
			panic(err)
		}

		return lst, nil
	}
}
