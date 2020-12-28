package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Uniform struct {
  typeExpr *TypeExpression
  nameExpr *VarExpression
  n int // -1 for scalar
  TokenData
}

func NewUniform(typeExpr *TypeExpression, name string, n int, ctx context.Context) *Uniform {
  return &Uniform{typeExpr, NewVarExpression(name, ctx), n, newTokenData(ctx)}
}

func (t *Uniform) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Uniform(")

  b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))

  if (t.n > 0) {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.n))
    b.WriteString("]")
  }

  b.WriteString(")\n")

  return b.String()
}

func (t *Uniform) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  // TODO: check if actually used
  b.WriteString(indent)
  b.WriteString("uniform ")
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())

  if (t.n > 0) {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.n))
    b.WriteString("]")
  }

  return b.String()
}
