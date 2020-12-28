package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type VarStatement struct {
  typeExpr *TypeExpression
  nameExpr *VarExpression
  rhsExpr Expression // optional, can be nil
  TokenData
}

func NewVarStatement(typeExpr *TypeExpression, name string, rhsExpr Expression, ctx context.Context) *VarStatement {
  return &VarStatement{typeExpr, NewVarExpression(name, ctx), rhsExpr, newTokenData(ctx)}
}

func (t *VarStatement) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("VarStatement(")
	b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))
	b.WriteString(")\n")

  if t.rhsExpr != nil {
		b.WriteString(t.rhsExpr.Dump(indent + "  "))
	}

	return b.String()
}

func (t *VarStatement) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(t.typeExpr.WriteExpression())
	b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())

  if t.rhsExpr != nil {
    b.WriteString("=")
    b.WriteString(t.rhsExpr.WriteExpression())
  }

	return b.String()
}
