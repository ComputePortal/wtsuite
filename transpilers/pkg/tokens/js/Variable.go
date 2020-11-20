package js

import (
  "./values"

	"../context"
)

// a Package also implements the Variable interface
type Variable interface {
	Context() context.Context
	Dump(indent string) string

	Name() string
	Rename(newName string)

	Constant() bool
	SetConstant()

  GetValue() values.Value
  SetValue(values.Value)

  // anything that can be evaluated during the resolve names stage (eg. class statement)
	GetObject() interface{} 
	SetObject(interface{})
}

type VariableData struct {
	name     string
	constant bool
  value    values.Value
	object   interface{}
	TokenData
}

func newVariableData(name string, constant bool, ctx context.Context) VariableData {
	return VariableData{name, constant, nil, nil, TokenData{ctx}}
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

// TODO: do this directly in the Namespace
func (t *VariableData) Rename(newName string) {
	t.name = newName
}

func (t *VariableData) Constant() bool {
	return t.constant
}

func (t *VariableData) SetConstant() {
	t.constant = true
}

func (t *VariableData) GetValue() values.Value {
  if t.value == nil {
    hereCtx := t.Context()
    here := hereCtx.NewError("here: " + t.Name())
    panic(here.Error())
  }

  return t.value
}

func (t *VariableData) SetValue(v values.Value) {
  t.value = v
}

func (t *VariableData) GetObject() interface{} {
	return t.object
}

func (t *VariableData) SetObject(ptr interface{}) {
	t.object = ptr
}
