package js

import (
	"strings"

	"./values"

	"../context"
)

type Call struct {
	lhs  Expression
	args []Expression
	TokenData
}

func NewCall(lhs Expression, args []Expression, ctx context.Context) *Call {
	return &Call{lhs, args, TokenData{ctx}}
}

// returns empty string if lhs is not *VarExpression
func (t *Call) Name() string {
	if ve, ok := t.lhs.(*VarExpression); ok {
		return ve.Name()
	} else {
		return ""
	}
}

func (t *Call) Args() []Expression {
	return t.args
}

func (t *Call) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Call\n")

	b.WriteString(t.lhs.Dump(indent + "  "))

	for _, arg := range t.args {
		b.WriteString(arg.Dump(indent + "( "))
	}

	return b.String()
}

func (t *Call) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.lhs.WriteExpression())

	b.WriteString("(")

	for i, arg := range t.args {
		b.WriteString(arg.WriteExpression())

		if i < len(t.args)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString(")")

	return b.String()
}

func (t *Call) WriteStatement(indent string) string {
	return indent + t.WriteExpression()
}

func (t *Call) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Call) HoistNames(scope Scope) error {
	return nil
}

func (t *Call) ResolveExpressionNames(scope Scope) error {
	// because sometimes both are possible, new needs to be available as well (for speed)
	if err := t.lhs.ResolveExpressionNames(scope); err != nil {
		return err
	}

	for _, arg := range t.args {
		if err := arg.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *Call) ResolveStatementNames(scope Scope) error {
	return t.ResolveExpressionNames(scope)
}

func (t *Call) evalArgs() ([]values.Value, error) {
	result := make([]values.Value, 0)

	for _, a := range t.args {
		val, err := a.EvalExpression()
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

func (t *Call) EvalExpression() (values.Value, error) {
	lhsVal, err := t.lhs.EvalExpression()
	if err != nil {
		return nil, err
	}

	args, err := t.evalArgs()
	if err != nil {
		return nil, err
	}

  res, err := lhsVal.EvalFunction(args, false, t.Context())
  if err != nil {
		if VERBOSITY >= 1 {
			context.AppendContextString(err, "Info: called here", t.Context())
		}

    return nil, err
  } else if res == nil {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: function returns void")
  }

  return res, nil
}

func (t *Call) EvalStatement() error {
	lhsVal, err := t.lhs.EvalExpression()
	if err != nil {
		return err
	}

	args, err := t.evalArgs()
	if err != nil {
		return err
	}

  res, err := lhsVal.EvalFunction(args, true, t.Context())
  if err != nil {
		if VERBOSITY >= 1 {
			context.AppendContextString(err, "Info: called here", t.Context())
		}

    return err
  } else if res != nil {
    errCtx := t.Context()
    return errCtx.NewError("Error: method returns non-void")
  }

	return nil
}

func (t *Call) ResolveExpressionActivity(usage Usage) error {
	if err := t.lhs.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	for _, arg := range t.args {
		if err := arg.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (t *Call) ResolveStatementActivity(usage Usage) error {
	return t.ResolveExpressionActivity(usage)
}

func (t *Call) UniversalExpressionNames(ns Namespace) error {
	if err := t.lhs.UniversalExpressionNames(ns); err != nil {
		return err
	}

	for _, arg := range t.args {
		if err := arg.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Call) UniqueExpressionNames(ns Namespace) error {
	if err := t.lhs.UniqueExpressionNames(ns); err != nil {
		return err
	}

	for _, arg := range t.args {
		if err := arg.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Call) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *Call) UniqueStatementNames(ns Namespace) error {
	return t.UniqueExpressionNames(ns)
}

func (t *Call) Walk(fn WalkFunc) error {
  if err := t.lhs.Walk(fn); err != nil {
    return err
  }

  for _, arg := range t.args {
    if err := arg.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}
