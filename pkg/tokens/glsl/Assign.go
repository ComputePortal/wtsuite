package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Assign struct {
  lhs Expression
  rhs Expression
  op string // eg. "+" for "+="
  TokenData
}

func NewAssign(lhs Expression, rhs Expression, op string, ctx context.Context) *Assign {
	if op == ":" || op == "!" || op == "=" || op == "==" || op == "!=" || op == ">" || op == "<" {
		err := ctx.NewError("not a valid assign op '" + op + "'")
		panic(err)
	}

  return &Assign{lhs, rhs, op, newTokenData(ctx)}
}

func (t *Assign) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Assign (")
	b.WriteString(t.op)
	b.WriteString("=\n")

  b.WriteString(t.lhs.Dump(indent + "  lhs:"))
  b.WriteString(t.rhs.Dump(indent + "  rhs:"))

	return b.String()
}

func (t *Assign) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

  b.WriteString(indent)
	b.WriteString(t.lhs.WriteExpression())
	b.WriteString(t.op)
	b.WriteString("=")
	b.WriteString(t.rhs.WriteExpression())

	return b.String()
}
