package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Parens struct {
  expr Expression
  TokenData
}

func NewParens(expr Expression, ctx context.Context) *Parens {
  return &Parens{expr, newTokenData(ctx)}
}

func (t *Parens) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Parens\n")

	b.WriteString(t.expr.Dump(indent + "(  "))

	return b.String()
}

func (t *Parens) WriteExpression() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(t.expr.WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (t *Parens) ResolveExpressionNames(scope Scope) error {
	return t.expr.ResolveExpressionNames(scope)
}

func (t *Parens) EvalExpression() (values.Value, error) {
  return t.expr.EvalExpression()
}

func (t *Parens) ResolveExpressionActivity(usage Usage) error {
  return t.expr.ResolveExpressionActivity(usage)
}
