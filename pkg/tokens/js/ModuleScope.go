package js

type ModuleScope struct {
	module  *ModuleData
	globals GlobalScope // both as parent and as globals
	ScopeData
}

func (ms *ModuleScope) HasVariable(name string) bool {
  res := ms.ScopeData.HasVariable(name)

  return res
}

func (ms *ModuleScope) GetVariable(name string) (Variable, error) {
  // injection of variabled handle by Import statement (so imported variables are correctly scoped)
  return ms.ScopeData.GetVariable(name)
}

func GetModule(s_ Scope) *ModuleData {
  if s, ok := s_.(*ModuleScope); ok {
    return s.module
  } else {
    if p := s.Parent(); p != nil {
      return GetModule(p)
    } else {
      return nil
    }
  }
}

func GetGlobalScope(s_ Scope) GlobalScope {
  if s, ok := s_.(*ModuleScope); ok {
    return s.globals
  } else {
    if p := s.Parent(); p != nil {
      return GetGlobalScope(p)
    } else {
      return nil
    }
  }
}
