package glsl

import (
	"strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type VarStatement struct {
  typeExpr *TypeExpression
  nameExpr *VarExpression
  length int
  rhsExpr Expression // optional, can be nil
  TokenData
}

func NewVarStatement(typeExpr *TypeExpression, name string, length int, rhsExpr Expression, ctx context.Context) *VarStatement {
  return &VarStatement{typeExpr, NewVarExpression(name, ctx), length, rhsExpr, newTokenData(ctx)}
}

func (t *VarStatement) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("VarStatement(")
	b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))

  if t.length > 0 {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

	b.WriteString(")\n")

  if t.rhsExpr != nil {
		b.WriteString(t.rhsExpr.Dump(indent + "  "))
	}

	return b.String()
}

func (t *VarStatement) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(t.typeExpr.WriteExpression())
	b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())

  if t.length > 0 {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  if t.rhsExpr != nil {
    b.WriteString("=")
    b.WriteString(t.rhsExpr.WriteExpression())
  }

	return b.String()
}

func (t *VarStatement) ResolveStatementNames(scope Scope) error {
  if err := t.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  if t.rhsExpr != nil {
    if err := t.rhsExpr.ResolveExpressionNames(scope); err != nil {
      return err
    }
  }

  variable := t.nameExpr.GetVariable()

  if err := scope.SetVariable(t.nameExpr.Name(), variable); err != nil {
    return err
  }

  return nil
}

func (t *VarStatement) EvalStatement() error {
  val, err := t.typeExpr.Instantiate(t.Context())
  if err != nil {
    return err
  }
  
  if t.length > 0 {
    val = values.NewArray(val, t.length, t.Context())
  }

  if t.rhsExpr != nil {
    rhsValue, err := t.rhsExpr.EvalExpression()
    if err != nil {
      return err
    }

    if err := val.Check(rhsValue, t.Context()); err != nil {
      return err
    }
  } 

  variable := t.nameExpr.GetVariable()
  variable.SetValue(val)

  return nil
}

func (t *VarStatement) ResolveStatementActivity(usage Usage) error {
  if t.rhsExpr != nil {
    if err := t.rhsExpr.ResolveExpressionActivity(usage); err != nil {
      return err
    }
  }

  if err := t.typeExpr.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  variable := t.nameExpr.GetVariable()

  return usage.Rereference(variable, t.Context())
}

func (t *VarStatement) UniqueStatementNames(ns Namespace) error {
  ns.VarName(t.nameExpr.GetVariable())

  return nil
}
