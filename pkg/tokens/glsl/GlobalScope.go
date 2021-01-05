package glsl

import ()

type GlobalScope interface {
  Scope
  GetModule(path string) (Module, error)
}

// wrapped by BundleScope in order to implement the GlobalScope interface
type GlobalScopeData struct {
	ScopeData
}
