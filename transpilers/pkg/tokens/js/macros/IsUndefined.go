package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type IsUndefined struct {
	varNotFound bool // if true, then always true
	isLiteral   bool // if true, then always false
	Macro
}

func NewIsUndefined(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &IsUndefined{false, false, newMacro(args, ctx)}, nil
}

func (m *IsUndefined) Dump(indent string) string {
	return indent + "IsUndefined(...)"
}

func (m *IsUndefined) WriteExpression() string {
	if m.varNotFound {
		return "true"
	} else if m.isLiteral {
		return "false"
	}

	var b strings.Builder

	b.WriteString("((")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(") === undefined)")

	return b.String()
}

func (m *IsUndefined) ResolveExpressionNames(scope js.Scope) error {
	if len(m.args) != 1 {
		errCtx := m.Context()
		return errCtx.NewError("Error: expected 1 argument")
	}

	arg_ := m.args[0]
	if arg, ok := arg_.(*js.VarExpression); ok {
		name := arg.Name()
		if !scope.HasVariable(name) {
			m.varNotFound = true
			return nil
		}
	}

	return arg_.ResolveExpressionNames(scope)
}

func (m *IsUndefined) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	if m.varNotFound {
		return prototypes.NewLiteralBoolean(true, ctx), nil
	}

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if values.IsLiteral(args[0]) {
		m.isLiteral = true
		return prototypes.NewLiteralBoolean(false, ctx), nil
	}

	return prototypes.NewBoolean(ctx), nil
}

func (m *IsUndefined) ResolveExpressionActivity(usage js.Usage) error {
	if m.varNotFound {
		return nil
	}

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *IsUndefined) UniversalExpressionNames(ns js.Namespace) error {
	if m.varNotFound {
		return nil
	}

	return m.Macro.UniversalExpressionNames(ns)
}

func (m *IsUndefined) UniqueExpressionNames(ns js.Namespace) error {
	if m.varNotFound {
		return nil
	}

	return m.Macro.UniqueExpressionNames(ns)
}