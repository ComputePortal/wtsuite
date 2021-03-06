package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Assign struct {
  lhs Expression
  rhs Expression
  op string // eg. "+" for "+="
  TokenData
}

func NewAssign(lhs Expression, rhs Expression, op string, ctx context.Context) *Assign {
	if op == ":" || op == "!" || op == "=" || op == "==" || op == "!=" || op == ">" || op == "<" {
		err := ctx.NewError("not a valid assign op '" + op + "'")
		panic(err)
	}

  return &Assign{lhs, rhs, op, newTokenData(ctx)}
}

func (t *Assign) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Assign (")
	b.WriteString(t.op)
	b.WriteString("=\n")

  b.WriteString(t.lhs.Dump(indent + "  lhs:"))
  b.WriteString(t.rhs.Dump(indent + "  rhs:"))

	return b.String()
}

func (t *Assign) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

  b.WriteString(indent)
	b.WriteString(t.lhs.WriteExpression())
	b.WriteString(t.op)
	b.WriteString("=")
	b.WriteString(t.rhs.WriteExpression())

	return b.String()
}

func (t *Assign) ResolveStatementNames(scope Scope) error {
	if err := t.lhs.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if err := t.rhs.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return nil
}

func (t *Assign) EvalStatement() error {
  rhsExpr := t.rhs

  if t.op != "" {
    switch t.op {
    case "+":
      rhsExpr = NewAddOp(t.lhs, t.rhs, t.Context())
    case "-":
      rhsExpr = NewSubOp(t.lhs, t.rhs, t.Context())
    default:
      errCtx := t.Context()
      return errCtx.NewError("Error: unrecognized assign op " + t.op)
    }
  }
  
  rhsValue, err := rhsExpr.EvalExpression()
  if err != nil {
    return err
  }

  switch lhsExpr := t.lhs.(type) {
  case *VarExpression:
    return lhsExpr.EvalSet(rhsValue, t.Context())
  case *Member:
    return lhsExpr.EvalSet(rhsValue, t.Context())
  case *Index:
    return lhsExpr.EvalSet(rhsValue, t.Context())
  default:
    errCtx := t.Context()
    return errCtx.NewError("Error: lhs not assignable")
  }
}

func (t *Assign) ResolveStatementActivity(usage Usage) error {
  if err := t.rhs.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  switch t.lhs.(type) {
  case *VarExpression:
    // nothing changes when VarExpression is lhs
  default:
    if err := t.lhs.ResolveExpressionActivity(usage); err != nil {
      return err
    }
  }

  return nil
}

func (t *Assign) UniqueStatementNames(ns Namespace) error {
  return nil
}
