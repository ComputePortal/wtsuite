package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type TypeExpression struct {
  VarExpression
}

func NewTypeExpression(name string, ctx context.Context) *TypeExpression {
  return &TypeExpression{newVarExpression(name, ctx)}
}

func (t *TypeExpression) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Type(")
  b.WriteString(t.VarExpression.Name())
  b.WriteString(")")

  return b.String()
}

func (t *TypeExpression) WriteExpression() string {
  return t.VarExpression.WriteExpression()
}
