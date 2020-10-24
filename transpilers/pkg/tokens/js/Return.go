package js

import (
	"strings"

	"./values"

	"../context"
)

type Return struct {
	expr Expression // can be nil for void return
	TokenData
}

func NewReturn(expr Expression, ctx context.Context) (*Return, error) {
	return &Return{expr, TokenData{ctx}}, nil
}

func (t *Return) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Return\n")

	if t.expr != nil {
		b.WriteString(t.expr.Dump(indent + "  "))
	}

	return b.String()
}

func (t *Return) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("return")

	if t.expr != nil {
		b.WriteString(" ")
		b.WriteString(t.expr.WriteExpression())
	}

	return b.String()
}

func (t *Return) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Return) HoistNames(scope Scope) error {
	return nil
}

func (t *Return) ResolveStatementNames(scope Scope) error {
	if t.expr != nil {
		return t.expr.ResolveExpressionNames(scope)
	}

	return nil
}

func IsVoidReturn(t Token) bool {
	if ret, ok := t.(*Return); ok {
		return ret.expr == nil
	}

	return false
}

func (t *Return) HoistValues(stack values.Stack) error {
	return nil
}

func (t *Return) EvalStatement(stack values.Stack) error {
	if t.expr == nil {
		retValue, err := stack.GetReturn(t.Context())
		if err != nil {
			return err
		}

		if retValue != nil && !values.IsVoid(retValue) {
			errCtx := t.Context()
			err := errCtx.NewError("Error: returning void, " +
				"in a function where previously non-void was returned")
			err.AppendContextString("Info: previous return value", retValue.Context())
			return err
		} else {
			return stack.SetReturn(values.NewVoid(t.Context()), t.Context())
		}
	}

	exprVal, err := t.expr.EvalExpression(stack)
	if err != nil {
		return err
	}

	retVal, err := stack.GetReturn(t.Context())
	if err != nil {
		return err
	}

	if retVal != nil && values.IsVoid(retVal) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: returning non-void, in a function where previously void was returned")
		err.AppendContextString("Info: previous return value", retVal.Context())
		return err
	}

	return stack.SetReturn(values.NewContextValue(exprVal, t.expr.Context()), t.Context())
}

func (t *Return) ResolveStatementActivity(usage Usage) error {
	if t.expr == nil {
		return nil
	}

	return t.expr.ResolveExpressionActivity(usage)
}

func (t *Return) UniversalStatementNames(ns Namespace) error {
	if t.expr != nil {
		if err := t.expr.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Return) UniqueStatementNames(ns Namespace) error {
	if t.expr != nil {
		if err := t.expr.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Return) Walk(fn WalkFunc) error {
  if t.expr != nil {
    if err := t.expr.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}
