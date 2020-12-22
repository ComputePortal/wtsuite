package js

import (
	"./prototypes"
)

type FunctionScope struct {
	function *Function
	ScopeData
}

func NewFunctionScope(fn *Function, parent Scope) *FunctionScope {
	return &FunctionScope{fn, newScopeData(parent)}
}

func (fs *FunctionScope) SetVariable(name string, v Variable) error {
	return fs.ScopeData.SetVariable(name, v)
}

func (fs *FunctionScope) IsAsync() bool {
	return prototypes.IsAsync(fs.function)
}

func (fs *FunctionScope) GetFunction() *Function {
  return fs.function
}
