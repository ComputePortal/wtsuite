package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type StatementScope struct {
	ScopeData
}

func newStatementScope(parent Scope) StatementScope {
	return StatementScope{newScopeData(parent)}
}

func NewStatementScope(parent Scope) *StatementScope {
	s := newStatementScope(parent)
	return &s
}

func (scope *StatementScope) GetNode() Node {
	return scope.parent.GetNode()
}

func (scope *StatementScope) Eval(key string, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return eval(scope, key, args, ctx)
}
