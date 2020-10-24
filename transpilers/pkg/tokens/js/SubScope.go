package js

type SubScope struct {
	ScopeData
}

// simply passes everyting to parent
func NewSubScope(parent Scope) *SubScope {
	return &SubScope{newScopeData(parent)}
}
