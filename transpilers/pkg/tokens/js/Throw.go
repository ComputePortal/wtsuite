package js

import (
	"strings"

	"./values"

	"../context"
)

type Throw struct {
	expr Expression
	TokenData
}

func NewThrow(expr Expression, ctx context.Context) (*Throw, error) {
	return &Throw{expr, TokenData{ctx}}, nil
}

func (t *Throw) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Throw\n")

	b.WriteString(t.expr.Dump(indent + "  "))

	return b.String()
}

func (t *Throw) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("throw ")
	b.WriteString(t.expr.WriteExpression())

	return b.String()
}

func (t *Throw) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Throw) HoistNames(scope Scope) error {
	return nil
}

func (t *Throw) ResolveStatementNames(scope Scope) error {
	return t.expr.ResolveExpressionNames(scope)
}

func (t *Throw) HoistValues(stack values.Stack) error {
	return nil
}

func (t *Throw) EvalStatement(stack values.Stack) error {
	exprValue, err := t.expr.EvalExpression(stack)
	if err != nil {
		return err
	}

	// TODO: should exprValue be instance of Error?

	if exprValue == nil {
		// should've been caught earlier
		errCtx := t.expr.Context()
		panic(errCtx.NewError("Error: expected non-void value"))
	}

	return nil
}

func (t *Throw) ResolveStatementActivity(usage Usage) error {
	return t.expr.ResolveExpressionActivity(usage)
}

func (t *Throw) UniversalStatementNames(ns Namespace) error {
	return t.expr.UniversalExpressionNames(ns)
}

func (t *Throw) UniqueStatementNames(ns Namespace) error {
	return t.expr.UniqueExpressionNames(ns)
}

func (t *Throw) Walk(fn WalkFunc) error {
  if err := t.expr.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}