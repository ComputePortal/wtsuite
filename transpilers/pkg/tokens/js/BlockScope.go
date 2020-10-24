package js

type BlockScope struct {
	ScopeData
}

func NewBlockScope(parent Scope) *BlockScope {
	return &BlockScope{newScopeData(parent)}
}

// vars are hoisted out (use SubScope to avoid this)
func (t *BlockScope) SetVariable(name string, v Variable) error {
	return t.Parent().SetVariable(name, v)
}
