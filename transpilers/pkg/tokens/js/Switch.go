package js

import (
	"strings"

	"./values"

	"../context"
)

type Switch struct {
	expr       Expression
	clauses    []Expression
	grouped    [][]Statement
	hasDefault bool
	Block      // dont use the Block.statements
}

func NewSwitch(expr Expression, ctx context.Context) (*Switch, error) {
	return &Switch{expr, make([]Expression, 0), make([][]Statement, 0), false, newBlock(ctx)}, nil
}

func (t *Switch) AddCase(clause Expression) error {
	if clause == nil {
		panic("nil not allowed")
	}

	t.clauses = append(t.clauses, clause)
	t.grouped = append(t.grouped, make([]Statement, 0))

	if len(t.clauses) != len(t.grouped) {
		panic("inconsistent lengths")
	}

	return nil
}

func (t *Switch) AddDefault(ctx context.Context) error {
	if t.hasDefault {
		return ctx.NewError("Error: default already defined")
	}

	t.hasDefault = true

	t.clauses = append(t.clauses, nil)
	t.grouped = append(t.grouped, make([]Statement, 0))

	return nil
}

func (t *Switch) AddStatement(statement Statement) {
	n := len(t.clauses)

	t.grouped[n-1] = append(t.grouped[n-1], statement)
}

func (t *Switch) ConvertToIf() *If {
	if t.expr != nil {
		panic("only convertible to If if there is no expr")
	}

	ifStatement, err := NewIf(t.Context())
	if err != nil {
		panic(err)
	}

	for i, group := range t.grouped {
		clause := t.clauses[i]
		if clause != nil {
			if err := ifStatement.AddCondition(clause); err != nil {
				panic("should've been caught before")
			}
		} else {
			if err := ifStatement.AddElse(); err != nil {
				panic("should've been caught before")
			}
		}

		for _, st := range group {
			ifStatement.AddStatement(st)
		}
	}

	return ifStatement
}

func (t *Switch) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Switch(")
	if t.expr != nil { // Switch might not have been converted to If yet, so allow expr==nil here
		b.WriteString("(")
		b.WriteString(strings.Replace(t.expr.Dump(""), "\n", "", -1))
		b.WriteString(")")
	}
	b.WriteString("\n")

	for i, clause := range t.clauses {
		b.WriteString(indent + "  ")
		if clause == nil {
			b.WriteString("Default\n")

			for _, statement := range t.grouped[i] {
				b.WriteString(statement.Dump(indent + "   :"))
			}
		} else {
			b.WriteString("Case(")
			b.WriteString(strings.Replace(clause.Dump(""), "\n", "", -1))
			b.WriteString(")\n")

			for _, statement := range t.grouped[i] {
				b.WriteString(statement.Dump(indent + "   :"))
			}
		}
	}

	return b.String()
}

func (t *Switch) WriteStatement(indent string) string {
	var b strings.Builder

	if t.expr == nil {
		panic("should've been converted to if")
	}

	b.WriteString(indent)
	b.WriteString("switch(")
	b.WriteString(t.expr.WriteExpression())
	b.WriteString("){")

	for i, clause := range t.clauses {
		b.WriteString(NL)
		b.WriteString(indent)
		b.WriteString(TAB)
		if clause == nil {
			b.WriteString("default")
		} else {
			b.WriteString("case ")
			b.WriteString(clause.WriteExpression())
		}

		if len(t.grouped[i]) > 0 {
			b.WriteString(":{")
			b.WriteString(NL)
			t.statements = t.grouped[i]
			b.WriteString(t.writeBlockStatements(indent+TAB+TAB, NL))
			b.WriteString(NL)
			b.WriteString(indent + TAB)
			b.WriteString("}")
		} else {
			b.WriteString(":")
		}
	}

	b.WriteString(NL)
	b.WriteString(indent)
	b.WriteString("}")

	return b.String()
}

func (t *Switch) HoistNames(scope Scope) error {
	for _, statements := range t.grouped {
		t.statements = statements
		if err := t.Block.HoistNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *Switch) ResolveStatementNames(scope Scope) error {
	if err := t.expr.ResolveExpressionNames(scope); err != nil {
		return err
	}

	for i, clause := range t.clauses {
		if clause != nil {
			if err := clause.ResolveExpressionNames(scope); err != nil {
				return err
			}
		}

		t.statements = t.grouped[i]

		subScope := NewCaseScope(scope)
		if err := t.Block.ResolveStatementNames(subScope); err != nil {
			return err
		}
	}

	return nil
}

func (t *Switch) EvalStatement() error {
	exprVal, err := t.expr.EvalExpression()
	if err != nil {
		return err
	}

  if !values.IsInstance(exprVal) {
		errCtx := exprVal.Context()
		return errCtx.NewError("Error: not a switchable value (" + exprVal.TypeName() + ")")
  }

	for i, clause := range t.clauses {
		if clause != nil {
			clauseVal, err := clause.EvalExpression()
			if err != nil {
				return err
			}

      if err := exprVal.Check(clauseVal, clauseVal.Context()); err != nil {
        return err
			}

			if err := t.Block.evalStatements(t.grouped[i]); err != nil {
				return err
			}
		} else if i < len(t.grouped) {
			if err := t.Block.evalStatements(t.grouped[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Switch) ResolveStatementActivity(usage Usage) error {
	for i := len(t.clauses) - 1; i >= 0; i-- {
		t.statements = t.grouped[i]
		if err := t.Block.ResolveStatementActivity(usage); err != nil {
			return err
		}

		clause := t.clauses[i]
		if clause != nil {
			if err := clause.ResolveExpressionActivity(usage); err != nil {
				return err
			}
		}
	}

	if err := t.expr.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *Switch) UniversalStatementNames(ns Namespace) error {
	if err := t.expr.UniversalExpressionNames(ns); err != nil {
		return err
	}

	for i, clause := range t.clauses {
		if clause != nil {
			if err := clause.UniversalExpressionNames(ns); err != nil {
				return err
			}
		}

		t.statements = t.grouped[i]
		if err := t.Block.UniversalStatementNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Switch) UniqueStatementNames(ns Namespace) error {
	if err := t.expr.UniqueExpressionNames(ns); err != nil {
		return err
	}

	for i, clause := range t.clauses {
		if clause != nil {
			if err := clause.UniqueExpressionNames(ns); err != nil {
				return err
			}
		}

		t.statements = t.grouped[i]

		subNs := ns.NewBlockNamespace()
		if err := t.Block.UniqueStatementNames(subNs); err != nil {
			return err
		}
	}

	return nil
}

func (t *Switch) Walk(fn WalkFunc) error {
  if err := t.expr.Walk(fn); err != nil {
    return err
  }

	for i, clause := range t.clauses {
		if clause != nil {
			if err := clause.Walk(fn); err != nil {
				return err
			}
		}

		t.statements = t.grouped[i]

		if err := t.Block.Walk(fn); err != nil {
			return err
		}
	}

	return fn(t)
}
