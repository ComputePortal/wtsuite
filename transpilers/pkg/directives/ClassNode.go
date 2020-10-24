package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
)

type ClassNode struct {
	parent     Node
	operations []Operation // XXX map with id might be better?
	thisEnums  map[string]*tokens.List
}

func NewClassNode(parent Node, operations []Operation) (*ClassNode, error) {
	return &ClassNode{parent, operations, make(map[string]*tokens.List)}, nil
}

func (n *ClassNode) Name() string {
	return n.parent.Name()
}

func (n *ClassNode) Type() NodeType {
	return n.parent.Type()
}

func (n *ClassNode) getElementCount() int {
	return n.parent.getElementCount()
}

func (n *ClassNode) getElementCountFolded() int {
	return n.parent.getElementCountFolded()
}

func (n *ClassNode) incrementElementCountFolded() {
	n.parent.incrementElementCountFolded()
}

func (n *ClassNode) getLastChild() tree.Tag {
	return n.parent.getLastChild()
}

func (n *ClassNode) AppendChild(child tree.Tag) error {
	return n.parent.AppendChild(child)
}

func (n *ClassNode) MapBlocks(blocks *tokens.StringDict) error {
	for _, op := range n.operations {
		if err := op.MapBlocks(blocks); err != nil {
			return err
		}
	}

	return n.parent.MapBlocks(blocks)
}

func (n *ClassNode) GetOperations() []Operation {
	return n.operations
}

func (n *ClassNode) PopOp(id string) (Operation, bool, error) {
	parentOp, parentOk, err := n.parent.PopOp(id)
	if err != nil {
		return nil, false, err
	}

	var thisOp Operation = nil
	thisOk := false
	for i, op := range n.operations {
		if op.ID() == id {
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

	if thisOk && parentOk {
		// merge
		mop, err := thisOp.Merge(parentOp)
		if err != nil {
			return nil, false, err
		}

		return mop, true, nil
	} else if thisOk {
		return thisOp, true, nil
	} else {
		return parentOp, parentOk, nil
	}
}

func (n *ClassNode) SearchStyle(key *tokens.String, ctx context.Context) (tokens.Token, error) {
	return n.parent.SearchStyle(key, ctx)
}

func (n *ClassNode) SearchAttrEnum(id *tokens.String, key *tokens.String, ctx context.Context) (*tokens.List, error) {
	if vals, ok := n.thisEnums[key.Value()]; ok {
		return vals, nil
	} else if id == nil {
		return nil, ctx.NewError("Error: attr enum " + key.Value() + " not found")
	}

	return n.parent.SearchAttrEnum(id, key, ctx)
}

func (n *ClassNode) addAttrEnum(key string, vals *tokens.List) {
	n.thisEnums[key] = vals
}
