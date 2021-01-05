package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
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

func (t *Call) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	return indent + t.WriteExpression()
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
	result := make([]values.Value, len(t.args))

	for i, a := range t.args {
		val, err := a.EvalExpression()
		if err != nil {
			return nil, err
		}

		result[i] = val
	}

	return result, nil
}

func (t *Call) eval() (values.Value, error) {
  fnVal, err := t.lhs.EvalExpression()
  if err != nil {
    return nil, err
  }

  argVals, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  return fnVal.EvalFunction(argVals, t.Context())
}

func (t *Call) EvalExpression() (values.Value, error) {
  retVal, err := t.eval()
  if err != nil {
    return nil, err
  }

  if retVal == nil {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: function returns void")
  }

  return retVal, nil
}

func (t *Call) EvalStatement() error {
  retVal, err := t.eval()
  if err != nil {
    return err
  }

  if retVal != nil {
    errCtx := t.Context()
    return errCtx.NewError("Error: function doesn't return void")
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

func (t *Call) UniqueStatementNames(ns Namespace) error {
  return nil
}
