package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Return struct {
  expr Expression // can be nil
  fn *Function
  TokenData
}

// expr can be nil
func NewReturn(expr Expression, ctx context.Context) *Return {
  return &Return{expr, nil, newTokenData(ctx)}
}

func (t *Return) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Return(")
  if t.expr != nil  {
    b.WriteString(t.expr.Dump(""))
  }
  b.WriteString(")")

  return b.String()
}

func (t *Return) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("return")
  if t.expr != nil {
    b.WriteString(" ")
    b.WriteString(t.expr.WriteExpression())
  }

  return b.String()
}

func (t *Return) ResolveStatementNames(scope Scope) error {
  fn := scope.GetFunction()
  if fn == nil {
    errCtx := t.Context()
    return errCtx.NewError("Error: return not inside function")
  }

  t.fn = fn

	if t.expr != nil {
    t.fn.RegisterReturn(t)
		return t.expr.ResolveExpressionNames(scope)
	}

	return nil
}

func (t *Return) EvalStatement() error {
  checkVal, err := t.fn.EvalCall(nil, t.Context())
  if err != nil {
    return err
  }

  if t.expr != nil {
    v, err := t.expr.EvalExpression()
    if err != nil {
      return err
    }

    if checkVal == nil {
      errCtx := t.Context()
      return errCtx.NewError("Error: function expects void return, got " + v.TypeName())
    } else if err := checkVal.Check(v, t.Context()); err != nil {
      return err
    } else {
      return nil
    }
  } else if checkVal != nil {
    errCtx := t.Context()
    return errCtx.NewError("Error: function expects " + checkVal.TypeName() + " return, got void")
  } else {
    return nil
  }
}

func (t *Return) ResolveStatementActivity(usage Usage) error {
  if t.expr == nil {
    return nil
  }

  return t.expr.ResolveExpressionActivity(usage)
}

func (t *Return) UniqueStatementNames(ns Namespace) error {
  return nil
}
