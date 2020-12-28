package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Varying struct {
  precType PrecisionType
  typeExpr *TypeExpression
  nameExpr *VarExpression
  TokenData
}

func NewVarying(precType PrecisionType, typeExpr *TypeExpression, name string, ctx context.Context) *Varying {
  return &Varying{precType, typeExpr, NewVarExpression(name, ctx), newTokenData(ctx)}
}

func (t *Varying) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Varying(")
  b.WriteString(PrecisionTypeToString(t.precType))
  b.WriteString(" ")
  b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))
  b.WriteString("\n")

  return b.String()
}

func (t *Varying) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  // TODO: check if actually used
  b.WriteString(indent)
  b.WriteString("varying ")
  if t.precType != DEFAULTP {
    b.WriteString(PrecisionTypeToString(t.precType))
    b.WriteString(" ")
  }
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())

  return b.String()
}
