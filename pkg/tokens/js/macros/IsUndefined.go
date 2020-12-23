package macros

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js"

	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type IsUndefined struct {
	varNotFound bool // if true, then always true
	isLiteral   bool // if true, then always false
	Macro
}

func NewIsUndefined(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

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

func (m *IsUndefined) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	if m.varNotFound {
		return prototypes.NewLiteralBoolean(true, ctx), nil
	}

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
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
