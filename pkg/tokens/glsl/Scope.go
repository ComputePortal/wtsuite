package glsl

type Scope interface {
  Parent() Scope

  HasVariable(name string) bool
  GetVariable(name string) (Variable, error)
  SetVariable(name string, v Variable) error

  GetFunction() *Function
}

type ScopeData struct {
  parent Scope
  variables map[string]Variable
}

func newScopeData(parent Scope) ScopeData {
	return ScopeData{parent, make(map[string]Variable)}
}

func NewScope(parent Scope) *ScopeData {
	scope := newScopeData(parent)
	return &scope
}

func (s *ScopeData) Parent() Scope {
	return s.parent
}

func (s *ScopeData) HasVariable(name string) bool {
	_, ok := s.variables[name]
	if !ok && s.parent != nil {
		return s.parent.HasVariable(name)
	} else {
		return ok
	}
}

func (s *ScopeData) GetVariable(name string) (Variable, error) {
	v, ok := s.variables[name]
	if !ok {
		if s.parent == nil {
			panic("should've been checked earlier")
		}

		return s.parent.GetVariable(name)
	}

	return v, nil // error can be used for import problems
}

func (s *ScopeData) SetVariable(name string, v Variable) error {
	if v == nil {
		panic("variable cant be nil")
	}

	if other, ok := s.variables[name]; ok {
		if other.Constant() {
			errCtx := v.Context()
			err := errCtx.NewError("Error: '" + name + "' previously set as const")
			err.AppendContextString("Info: defined here", other.Context())
			return err
		}
	}

	// variables can actually be overwritten (for eg. using 'let i' twice)
	s.variables[name] = v

	return nil
}

func (s *ScopeData) GetFunction() *Function {
  if s.parent != nil {
    return s.parent.GetFunction()
  } else {
    return nil
  }
}
