package js

import (
	"strings"

	"./values"

	"../context"
)

type Parens struct {
	expr Expression
	TokenData
}

func NewParens(expr Expression, ctx context.Context) *Parens {
	return &Parens{expr, TokenData{ctx}}
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

func (t *Parens) CollectTypeGuards(c map[Variable]values.Interface) (bool, error) {
	if expr, ok := t.expr.(TypeGuard); ok {
		return expr.CollectTypeGuards(c)
	} else {
		return false, nil
	}
}

func (t *Parens) ResolveExpressionActivity(usage Usage) error {
	return t.expr.ResolveExpressionActivity(usage)
}

func (t *Parens) UniversalExpressionNames(ns Namespace) error {
	return t.expr.UniversalExpressionNames(ns)
}

func (t *Parens) UniqueExpressionNames(ns Namespace) error {
	return t.expr.UniqueExpressionNames(ns)
}

func (t *Parens) Walk(fn WalkFunc) error {
  if err := t.expr.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
