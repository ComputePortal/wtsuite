package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type While struct {
	cond Expression
	Block
}

func NewWhile(cond Expression, ctx context.Context) (*While, error) {
	return &While{cond, newBlock(ctx)}, nil
}

func (t *While) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("While(")
	b.WriteString(strings.Replace(t.cond.WriteExpression(), "\n", "", -1))

	b.WriteString(")\n")

	for _, s := range t.statements {
		b.WriteString(s.Dump(indent + "{ "))
	}

	return b.String()
}

func (t *While) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("while(")
	b.WriteString(t.cond.WriteExpression())
	b.WriteString("){")
	b.WriteString(NL)

	b.WriteString(t.writeBlockStatements(indent+TAB, NL))

	b.WriteString(NL)
	b.WriteString(indent)
	b.WriteString("}")

	return b.String()
}

func (t *While) HoistNames(scope Scope) error {
	return t.Block.HoistNames(scope)
}

func (t *While) ResolveStatementNames(scope Scope) error {
	if err := t.cond.ResolveExpressionNames(scope); err != nil {
		return err
	}

	subScope := NewLoopScope(scope)

	return t.Block.ResolveStatementNames(subScope)
}

func (t *While) EvalStatement() error {
	condVal, err := t.cond.EvalExpression()
	if err != nil {
		return err
	}

	if !prototypes.IsBoolean(condVal) {
		errCtx := condVal.Context()
		return errCtx.NewError("Error: expected boolean condition")
	}

	return t.Block.EvalStatement()
}

func (t *While) ResolveStatementActivity(usage Usage) error {
	if err := t.Block.ResolveStatementActivity(usage); err != nil {
		return err
	}

	if err := t.cond.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *While) UniversalStatementNames(ns Namespace) error {
	if err := t.cond.UniversalExpressionNames(ns); err != nil {
		return err
	}

	return t.Block.UniversalStatementNames(ns)
}

func (t *While) UniqueStatementNames(ns Namespace) error {
	if err := t.cond.UniqueExpressionNames(ns); err != nil {
		return err
	}

	subNs := ns.NewBlockNamespace()

	return t.Block.UniqueStatementNames(subNs)
}

func (t *While) Walk(fn WalkFunc) error {
  if err := t.cond.Walk(fn); err != nil {
    return err
  }

  if err := t.Block.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
