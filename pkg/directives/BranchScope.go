package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
)

// BranchScope sends to exports to filescope (regular scopedata gives error)
type BranchScope struct {
  ScopeData
}

func NewBranchScope(parent Scope) Scope {
	if parent == nil {
		panic("parent can't be nil")
	}

	return &BranchScope{newScopeData(parent)}
}

func (s *BranchScope) SetVar(key string, v functions.Var) error {
  // this goes all the way to the filescope, so variables can be redefined
  if v.Exported {
    return s.parent.SetVar(key, v)
  } else {
    return s.ScopeData.SetVar(key, v)
  }
}

func (s *BranchScope) SetTemplate(key string, d Template) error {
  // this goes all the way to the filescope, so variables can be redefined
  if d.exported {
    return s.parent.SetTemplate(key, d)
  } else {
    return s.ScopeData.SetTemplate(key, d)
  }
}
