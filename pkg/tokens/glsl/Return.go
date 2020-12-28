package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Return struct {
  expr Expression // can be nil
  TokenData
}

// expr can be nil
func NewReturn(expr Expression, ctx context.Context) *Return {
  return &Return{expr, newTokenData(ctx)}
}

func (t *Return) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Return(")
  if t.expr != nil  {
    b.WriteString(t.expr.Dump(""))
  }
  b.WriteString(")")

  return b.String()
}

func (t *Return) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("return")
  if t.expr != nil {
    b.WriteString(" ")
    b.WriteString(t.expr.WriteExpression())
  }

  return b.String()
}
