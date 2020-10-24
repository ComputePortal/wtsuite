package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type If struct {
	conds   []Expression
	grouped [][]Statement
	Block   // dont use the Block.statements
}

func NewIf(ctx context.Context) (*If, error) {
	return &If{make([]Expression, 0), make([][]Statement, 0), newBlock(ctx)}, nil
}

func (t *If) AddCondition(expr Expression) error {
	if expr == nil {
		panic("nil not allowed")
	}

	t.conds = append(t.conds, expr)
	t.grouped = append(t.grouped, make([]Statement, 0))

	if len(t.conds) != len(t.grouped) {
		panic("inconsistent lengths")
	}

	return nil
}

func (t *If) AddElse() error {
	if t.conds[len(t.conds)-1] == nil {
		panic("else already added")
	}

	t.conds = append(t.conds, nil)
	t.grouped = append(t.grouped, make([]Statement, 0))

	return nil
}

func (t *If) AddStatement(statement Statement) {
	n := len(t.conds)

	t.grouped[n-1] = append(t.grouped[n-1], statement)
}

func (t *If) Dump(indent string) string {
	var b strings.Builder

	for i, c := range t.conds {
		b.WriteString(indent)
		if i == 0 {
			b.WriteString("If(")
			b.WriteString(strings.Replace(c.WriteExpression(), "\n", "", -1))
			b.WriteString(")\n")
		} else if c == nil {
			if i != len(t.conds)-1 {
				panic("only last can be nil")
			}
			b.WriteString("Else\n")
		} else {
			b.WriteString("ElseIf(")
			b.WriteString(strings.Replace(c.WriteExpression(), "\n", "", -1))
			b.WriteString(")\n")
		}

		for _, statement := range t.grouped[i] {
			b.WriteString(statement.Dump(indent + "{ "))

		}
	}

	return b.String()
}

func (t *If) WriteStatement(indent string) string {
	var b strings.Builder

	for i, c := range t.conds {
		if i == 0 {
			b.WriteString(indent)
			b.WriteString("if(")
			b.WriteString(c.WriteExpression())
			b.WriteString(")")
		} else if c != nil {
			b.WriteString(NL)
			b.WriteString(indent)
			b.WriteString("else if(")
			b.WriteString(c.WriteExpression())
			b.WriteString(")")
		} else {
			b.WriteString(NL)
			b.WriteString(indent)
			b.WriteString("else")
		}

		t.statements = t.grouped[i]
		b.WriteString("{")
		b.WriteString(NL)
		b.WriteString(t.writeBlockStatements(indent+TAB, NL))
		b.WriteString(NL)
		b.WriteString(indent)
		b.WriteString("}")
	}

	return b.String()
}

///////////////////////////
// 1. Name resolution stage
///////////////////////////
func (t *If) HoistNames(scope Scope) error {
	for _, statements := range t.grouped {
		t.statements = statements
		if err := t.Block.HoistNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *If) ResolveStatementNames(scope Scope) error {
	for i, cond := range t.conds {
		if cond != nil {
			if err := cond.ResolveExpressionNames(scope); err != nil {
				return err
			}
		}

		t.statements = t.grouped[i]

		subScope := NewBranchScope(scope)
		if err := t.Block.ResolveStatementNames(subScope); err != nil {
			return err
		}
	}

	return nil
}

func (t *If) HoistValues(stack values.Stack) error {
	return nil
}

func (t *If) evalTypeGuards(stack values.Stack, cond Expression) (map[interface{}]values.Interface, error) {
	if cond == nil {
		return nil, nil
	}

	if typeGuard, ok := cond.(TypeGuard); ok {
		typeGuards := make(map[interface{}]values.Interface)

		isTG, err := typeGuard.CollectTypeGuards(stack, typeGuards)
		if err != nil {
			return nil, err
		}

		if isTG {
			if len(typeGuards) == 0 {
				// some conditions can work alongside typeguards, even though they dont add any typeguards
				return nil, nil
			} else {
				return typeGuards, nil
			}
		} else {
			return nil, nil
		}
	} else {
		return nil, nil
	}
}

func (t *If) EvalStatement(stack values.Stack) error {
	for i, cond := range t.conds {
		condIsLit := false
		condLitVal := false
		typeGuards, err := t.evalTypeGuards(stack, cond)
		if err != nil {
			return err
		}

		if cond != nil && typeGuards == nil { // cond == nil -> else {...}
			// condition cannot be literal if there are typeGuards present
			condVal, err := cond.EvalExpression(stack)
			if err != nil {
				return err
			}

			if !condVal.IsInstanceOf(prototypes.Boolean) {
				errCtx := condVal.Context()
				return errCtx.NewError("Error: expected boolean condition")
			}

			condLitVal, condIsLit = condVal.LiteralBooleanValue()

			if condIsLit && !condLitVal {
				continue
			}
		}

		subStack := NewTypeGuardBranchStack(typeGuards, stack)

		if err := t.Block.evalStatements(t.grouped[i], subStack); err != nil {
			return err
		}

		if condIsLit && condLitVal {
			break
		}
	}

	return nil
}

func (t *If) ResolveStatementActivity(usage Usage) error {
	for i := len(t.conds) - 1; i >= 0; i-- {
		t.statements = t.grouped[i]
		if err := t.Block.ResolveStatementActivity(usage); err != nil {
			return err
		}

		cond := t.conds[i]
		if cond != nil {
			if err := cond.ResolveExpressionActivity(usage); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *If) UniversalStatementNames(ns Namespace) error {
	for i, cond := range t.conds {
		if cond != nil {
			if err := cond.UniversalExpressionNames(ns); err != nil {
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

func (t *If) UniqueStatementNames(ns Namespace) error {
	for i, cond := range t.conds {
		if cond != nil {
			if err := cond.UniqueExpressionNames(ns); err != nil {
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

func (t *If) Walk(fn WalkFunc) error {
  for i, cond := range t.conds {
    if cond != nil {
      if err := cond.Walk(fn); err != nil {
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
