package js

type CaseScope struct {
	ScopeData
}

func NewCaseScope(parent Scope) *CaseScope {
	return &CaseScope{newScopeData(parent)}
}

func (scope *CaseScope) IsBreakable() bool {
	return true
}

func (scope *CaseScope) IsContinueable() bool {
	return scope.parent.IsContinueable()
}

func (scope *CaseScope) IsAsync() bool {
	return scope.parent.IsAsync()
}
