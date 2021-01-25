package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Precision struct {
  precType PrecisionType
  typeExpr *TypeExpression
  TokenData
}

func NewPrecision(precType PrecisionType, typeExpr *TypeExpression, ctx context.Context) *Precision {
  return &Precision{precType, typeExpr, newTokenData(ctx)}
}

func (t *Precision) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Precision(")
  b.WriteString(PrecisionTypeToString(t.precType))
  b.WriteString(") ")
  b.WriteString(t.typeExpr.Dump(""))
  b.WriteString("\n")

  return b.String()
}

func (t *Precision) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("precision ")
  b.WriteString(PrecisionTypeToString(t.precType))
  b.WriteString(" ")
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(";")

  return b.String()
}

func (t *Precision) ResolveStatementNames(scope Scope) error {
	return t.typeExpr.ResolveExpressionNames(scope)
}

func (t *Precision) EvalStatement() error {
  return nil
}

func (t *Precision) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *Precision) UniqueStatementNames(ns Namespace) error {
  return nil
}
