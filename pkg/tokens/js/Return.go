package js

import (
	"strings"

  "./values"

	"../context"
)

type Return struct {
	expr Expression // can be nil for void return
  fn   *Function // registered during resolve stage
	TokenData
}

func NewReturn(expr Expression, ctx context.Context) (*Return, error) {
	return &Return{expr, nil, TokenData{ctx}}, nil
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

func IsVoidReturn(t Token) bool {
	if ret, ok := t.(*Return); ok {
		return ret.expr == nil
	}

	return false
}

func (t *Return) EvalStatement() error {
  var exprVal values.Value = nil
	if t.expr != nil {
    var err error
    exprVal, err = t.expr.EvalExpression()
    if err != nil {
      return err
    }
  }

  if t.fn == nil {
    hereCtx := t.Context()
    here := hereCtx.NewError("fn should've been found")
    panic(here.Error())
  }

  retVal, err := t.fn.getReturnValue()
  if err != nil {
    return err
  }

  if retVal == nil {
    if exprVal != nil {
      errCtx := t.Context()
      return errCtx.NewError("Error: expected void return value")
    }
  } else {
    if exprVal == nil {
      errCtx := t.Context()
      return errCtx.NewError("Error: unexpected return value")
    }

    if err := retVal.Check(exprVal, t.Context()); err != nil {
      return err
    }
  }

  return nil
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
