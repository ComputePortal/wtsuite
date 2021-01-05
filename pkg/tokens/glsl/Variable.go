package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Variable interface {
  Context() context.Context
  Name() string
  Rename(string)

  Constant() bool
  SetConstant()

  GetValue() values.Value
  SetValue(val values.Value)

  GetObject() interface{}
  SetObject(obj interface{})
}

type VariableData struct {
  name string
  constant bool
  value values.Value
  obj interface{}
  TokenData
}

func newVariableData(name string, constant bool, ctx context.Context) VariableData {
  return VariableData{name, constant, nil, nil, newTokenData(ctx)}
}

func NewVariable(name string, ctx context.Context) *VariableData {
  res := newVariableData(name, false, ctx)

  return &res
}

func (v *VariableData) Name() string {
  return v.name
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
    here := hereCtx.NewError("Internal Error: " + t.Name() + " value not set")
    panic(here.Error())
  }

  return t.value
}

func (t *VariableData) SetValue(v values.Value) {
  t.value = v
}

func (t *VariableData) Rename(newName string) {
  t.name = newName
}

func (t *VariableData) SetObject(obj interface{}) {
  t.obj = obj
}

func (t *VariableData) GetObject() interface{} {
  return t.obj
}
