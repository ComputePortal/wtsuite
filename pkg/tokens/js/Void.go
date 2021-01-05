package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Void struct {
	expr Expression // can be nil for void return
	TokenData
}

func NewVoidStatement(expr Expression, ctx context.Context) *Void {
	return &Void{expr, TokenData{ctx}}
}

func (t *Void) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Void\n")

	b.WriteString(t.expr.Dump(indent + "  "))

	return b.String()
}

func (t *Void) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	// it is unnecesary to write "void" in the final output
	b.WriteString(t.expr.WriteExpression())

	return b.String()
}

func (t *Void) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Void) HoistNames(scope Scope) error {
	return nil
}

func (t *Void) ResolveStatementNames(scope Scope) error {
	return t.expr.ResolveExpressionNames(scope)
}

func (t *Void) EvalStatement() error {
	v, err := t.expr.EvalExpression()

	if err != nil {
		return err
	}

	// void expects a non-nil return value
	if v == nil {
		errCtx := t.Context()
		return errCtx.NewError("Error: void argument doesn't return a value")
	}

	return nil
}

func (t *Void) ResolveStatementActivity(usage Usage) error {
	return t.expr.ResolveExpressionActivity(usage)
}

func (t *Void) UniversalStatementNames(ns Namespace) error {
	return t.expr.UniversalExpressionNames(ns)
}

func (t *Void) UniqueStatementNames(ns Namespace) error {
	return t.expr.UniqueExpressionNames(ns)
}

func (t *Void) Walk(fn WalkFunc) error {
  if err := t.expr.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
