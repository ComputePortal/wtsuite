package js

type LoopScope struct {
	ScopeData
}

func NewLoopScope(parent Scope) *LoopScope {
	return &LoopScope{newScopeData(parent)}
}

func (scope *LoopScope) IsBreakable() bool {
	return true
}

func (scope *LoopScope) IsContinueable() bool {
	return true
}

func (scope *LoopScope) IsAsync() bool {
	return scope.parent.IsAsync()
}
