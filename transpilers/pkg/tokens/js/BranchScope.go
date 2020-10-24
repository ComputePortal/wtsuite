package js

type BranchScope struct {
	ScopeData
}

func NewBranchScope(parent Scope) *BranchScope {
	return &BranchScope{newScopeData(parent)}
}

func (scope *BranchScope) IsBreakable() bool {
	return scope.parent.IsBreakable()
}

func (scope *BranchScope) IsContinueable() bool {
	return scope.parent.IsContinueable()
}

func (scope *BranchScope) IsAsync() bool {
	return scope.parent.IsAsync()
}
