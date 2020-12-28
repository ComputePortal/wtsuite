package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Variable interface {
  Context() context.Context
  Name() string
}

type VariableData struct {
  name string
  TokenData
}

func newVariableData(name string, ctx context.Context) VariableData {
  return VariableData{name, newTokenData(ctx)}
}

func NewVariable(name string, ctx context.Context) *VariableData {
  res := newVariableData(name, ctx)

  return &res
}

func (v *VariableData) Name() string {
  return v.name
}
