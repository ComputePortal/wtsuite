package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type ClassScope struct {
	class *Class
	ScopeData
}

func NewClassScope(class *Class, parent Scope) *ClassScope {
	return &ClassScope{class, newScopeData(parent)}
}

func (s *ClassScope) FriendlyPrototypes() []values.Prototype {
	res := s.ScopeData.FriendlyPrototypes()

	return append(res, s.class)
}
