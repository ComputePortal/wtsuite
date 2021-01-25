package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Attribute struct {
  Pointer
}

func NewAttribute(typeExpr *TypeExpression, name string, ctx context.Context) *Attribute {
  return &Attribute{newPointer(typeExpr, NewVarExpression(name, ctx), -1, ctx)}
}

func (t *Attribute) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Attribute")
  b.WriteString("\n")
  b.WriteString(t.typeExpr.Dump(indent + "  "))
  b.WriteString("\n")
  b.WriteString(t.nameExpr.Dump(indent+"  "))

  return b.String()
}

func (t *Attribute) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  // TODO: check if actually used
  b.WriteString(indent)
  b.WriteString("attribute ")
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())
  b.WriteString(";")

  return b.String()
}
