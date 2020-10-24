package js

import (
	"../context"
)

// a Package also implements the Variable interface
type Variable interface {
	Context() context.Context
	Dump(indent string) string
	Name() string
	Constant() bool
	SetConstant()
	Rename(newName string)
	GetObject() interface{} // anything that can be evaluated during the resolve names stage (eg. class statement)
	SetObject(interface{})
}

type VariableData struct {
	name     string
	constant bool
	object   interface{}

	TokenData
}

func newVariableData(name string, constant bool, ctx context.Context) VariableData {
	return VariableData{name, constant, nil, TokenData{ctx}}
}

func NewVariable(name string, constant bool, ctx context.Context) *VariableData {
	res := newVariableData(name, constant, ctx)
	return &res
}

func (t *VariableData) Dump(indent string) string {
	return indent + "Variable " + t.name
}

func (t *VariableData) Name() string {
	return t.name
}

func (t *VariableData) Constant() bool {
	return t.constant
}

func (t *VariableData) SetConstant() {
	t.constant = true
}

// TODO: do this directly in the Namespace
func (t *VariableData) Rename(newName string) {
	t.name = newName
}

func (t *VariableData) GetObject() interface{} {
	return t.object
}

func (t *VariableData) SetObject(ptr interface{}) {
	t.object = ptr
}
