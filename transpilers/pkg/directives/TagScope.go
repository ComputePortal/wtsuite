package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type TagScope struct {
	node Node
	ScopeData
}

func newTagScope(parent Scope, node Node) *TagScope {
	return &TagScope{node, newScopeData(parent)}
}

func NewSubScope(parent Scope, node Node) *TagScope {
	if parent == nil {
		panic("parent can't be nil")
	}

	if node == nil {
		panic("node can't be nil")
	}

	return newTagScope(parent, node)
}

func NewRootScope(node *RootNode) *TagScope {
	return newTagScope(nil, node)
}

func (scope *TagScope) GetNode() Node {
	return scope.node
}

func (scope *TagScope) Eval(key string, args []tokens.Token,
	ctx context.Context) (tokens.Token, error) {
	return eval(scope, key, args, ctx)
}
