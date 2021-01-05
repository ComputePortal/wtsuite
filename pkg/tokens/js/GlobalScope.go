package js

import ()

type GlobalScope interface {
	Scope
	GetModule(path string) (Module, error)
}

var ActivateMacroHeaders func(name string) = nil

// wrapped by BundleScope in order to implement the GlobalScope interface
type GlobalScopeData struct {
	ScopeData
}

func (s *GlobalScopeData) GetVariable(name string) (Variable, error) {
	ActivateMacroHeaders(name)

	return s.ScopeData.GetVariable(name)
}
