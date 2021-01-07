package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type If struct {
  conds []Expression
  groups []*Block
  TokenData
}

func NewIf(ctx context.Context) *If {
  return &If{make([]Expression, 0), make([]*Block, 0), newTokenData(ctx)}
}

func (t *If) AddCondition(expr Expression) error {
	if expr == nil {
		panic("nil not allowed")
	}

	t.conds = append(t.conds, expr)
	t.groups = append(t.groups, NewBlock(t.Context()))

	if len(t.conds) != len(t.groups) {
		panic("inconsistent lengths")
	}

	return nil
}

func (t *If) AddElse() error {
  if t.conds[len(t.conds)-1] == nil {
    panic("else already added")
  }

	t.conds = append(t.conds, nil)
	t.groups = append(t.groups, NewBlock(t.Context()))

	return nil
}

func (t *If) AddStatement(statement Statement) {
	n := len(t.conds)

	t.groups[n-1].AddStatement(statement)
}

func (t *If) Dump(indent string) string {
	var b strings.Builder

	for i, c := range t.conds {
		b.WriteString(indent)
		if i == 0 {
			b.WriteString("If(")
			b.WriteString(strings.Replace(c.Dump(""), "\n", "", -1))
			b.WriteString(")\n")
		} else if c == nil {
			if i != len(t.conds)-1 {
				panic("only last can be nil")
			}
			b.WriteString("Else\n")
		} else {
			b.WriteString("ElseIf(")
			b.WriteString(strings.Replace(c.Dump(""), "\n", "", -1))
			b.WriteString(")\n")
		}

    b.WriteString(t.groups[i].Dump(indent + "{ "))
	}

	return b.String()
}

func (t *If) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	for i, c := range t.conds {
		if i == 0 {
			b.WriteString(indent)
			b.WriteString("if(")
			b.WriteString(c.WriteExpression())
			b.WriteString(")")
		} else if c != nil {
			b.WriteString(nl)
			b.WriteString(indent)
			b.WriteString("else if(")
			b.WriteString(c.WriteExpression())
			b.WriteString(")")
		} else {
			b.WriteString(nl)
			b.WriteString(indent)
			b.WriteString("else")
		}

		b.WriteString("{")
		b.WriteString(nl)
		b.WriteString(t.groups[i].writeBlockStatements(usage, indent+tab, nl, tab))
		b.WriteString(nl)
		b.WriteString(indent)
		b.WriteString("}")
	}

	return b.String()
}

func (t *If) ResolveStatementNames(scope Scope) error {
	for i, cond := range t.conds {
		if cond != nil {
			if err := cond.ResolveExpressionNames(scope); err != nil {
				return err
			}
		}

    block := t.groups[i]

		subScope := NewScope(scope)
		if err := block.ResolveStatementNames(subScope); err != nil {
			return err
		}
	}

	return nil
}

func (t *If) EvalStatement() error {
  for i, cond := range t.conds {
    if cond != nil {
      condVal, err := cond.EvalExpression()
      if err != nil {
        return err
      }

      if !values.IsBool(condVal) {
        errCtx := cond.Context()
        return errCtx.NewError("Error: expected bool condition, got " + condVal.TypeName())
      }
    }

    if err := t.groups[i].evalStatements(); err != nil {
      return err
    }
  }

  return nil
}

func (t *If) ResolveStatementActivity(usage Usage) error {
	for i := len(t.conds) - 1; i >= 0; i-- {
		if err := t.groups[i].ResolveStatementActivity(usage); err != nil {
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

func (t *If) UniqueStatementNames(ns Namespace) error {
	for _, group := range t.groups {
		subNs := ns.NewBlockNamespace()
		if err := group.UniqueStatementNames(subNs); err != nil {
			return err
		}
	}

	return nil
}
