package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// common base class for ForIn and ForOf
type ForInOf struct {
	lhs *VarExpression
	rhs Expression
	ForBlock
}

func newForInOf(varType VarType, lhs *VarExpression, rhs Expression, ctx context.Context) ForInOf {
	return ForInOf{lhs, rhs, newForBlock(varType, ctx)}
}

func (t *ForInOf) dump(indent string, op string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("For ")
	b.WriteString(VarTypeToString(t.varType))
	b.WriteString(" ")
	b.WriteString(op)
	b.WriteString("\n")

	b.WriteString(strings.Replace(t.rhs.Dump(""), "\n", "", -1))

	b.WriteString("\n")

	for _, s := range t.statements {
		b.WriteString(s.Dump(indent + "{ "))
	}

	return b.String()
}

func (t *ForInOf) writeStatement(usage Usage, indent string, extra string, op string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(t.writeStatementHeader(indent, extra, true))

	b.WriteString(t.lhs.WriteExpression())
	b.WriteString(" ")
	b.WriteString(op)
	b.WriteString(" ")
	b.WriteString(t.rhs.WriteExpression())
	b.WriteString(t.writeStatementFooter(usage, indent, nl, tab))

	return b.String()
}

///////////////////////////
// 1. Name resolution stage
///////////////////////////
func (t *ForInOf) HoistNames(scope Scope) error {
	if t.varType == VAR {
		if err := scope.SetVariable(t.lhs.Name(), t.lhs.GetVariable()); err != nil {
			return err
		}
	}

	return t.Block.HoistNames(scope)
}

func (t *ForInOf) ResolveStatementNames(scope Scope) error {
	subScope := NewLoopScope(scope)

	name := t.lhs.Name()

	switch t.varType {
	case LET, CONST:
		if err := subScope.SetVariable(name, t.lhs.GetVariable()); err != nil {
			return err
		}
	case VAR:
		if !scope.HasVariable(name) {
			panic("should've been hoisted before")
		}
	default:
		panic("unhandled")
	}

	if err := t.rhs.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return t.Block.ResolveStatementNames(subScope)
}

func (t *ForInOf) ResolveStatementActivity(usage Usage) error {
	if err := t.Block.ResolveStatementActivity(usage); err != nil {
		return err
	}

	if err := usage.Rereference(t.lhs.GetVariable(), t.lhs.Context()); err != nil {
		return err
	}

	if err := t.rhs.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *ForInOf) UniversalStatementNames(ns Namespace) error {
	if err := t.lhs.UniversalExpressionNames(ns); err != nil {
		return err
	}

	if err := t.rhs.UniversalExpressionNames(ns); err != nil {
		return err
	}

	return t.Block.UniversalStatementNames(ns)
}

func (t *ForInOf) UniqueStatementNames(ns Namespace) error {
	subNs := ns.NewBlockNamespace()

	ref := t.lhs.GetVariable()

	switch t.varType {
	case LET, CONST:
		subNs.LetName(ref)
	case VAR:
		ns.VarName(ref)
	default:
		panic("unexpected")
	}

	if err := t.rhs.UniqueExpressionNames(ns); err != nil {
		return err
	}

	return t.Block.UniqueStatementNames(subNs)
}

func (t *ForInOf) Walk(fn WalkFunc) error {
  if err := t.lhs.Walk(fn); err != nil {
    return err
  }

  if err := t.rhs.Walk(fn); err != nil {
    return err
  }

  return t.ForBlock.Walk(fn)
}
