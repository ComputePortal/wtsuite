package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type VarExpression struct {
  variable Variable
  origName string

  TokenData
}

func newVarExpression(name string, ctx context.Context) VarExpression {
  return VarExpression{NewVariable(name, ctx), name, newTokenData(ctx)}
}

func NewVarExpression(name string, ctx context.Context) *VarExpression {
  ve := newVarExpression(name, ctx)

  return &ve
}

func (t *VarExpression) Dump(indent string) string {
  s := indent + "Var(" + t.Name() + ")\n"
  return s
}

func (t *VarExpression) GetVariable() Variable {
  return t.variable
}

func (t *VarExpression) Name() string {
  return t.variable.Name()
}

func (t *VarExpression) WriteExpression() string {
  return t.Name()
}
