package macros

import (
	"fmt"

	"../"

	"../values"

	"../../context"
)

type StackTrace struct {
	Macro
}

func NewStackTrace(args []js.Expression, ctx context.Context) (js.Expression, error) {
	fmt.Println("creating new stacktrace")
	return &StackTrace{newMacro(args, ctx)}, nil
}

func (m *StackTrace) Dump(indent string) string {
	return indent + "StackTrace(...)"
}

func (m *StackTrace) WriteExpression() string {
	return ""
}

func (m *StackTrace) ResolveExpressionNames(scope js.Scope) error {
	if len(m.args) != 0 {
		errCtx := m.Context()
		return errCtx.NewError("Error: expected 0 arguments")
	}

	return m.Macro.ResolveExpressionNames(scope)
}

func (m *StackTrace) EvalExpression(stack values.Stack) (values.Value, error) {
	errCtx := m.Context()
	return nil, errCtx.NewError("Info: stack start")
}

func (m *StackTrace) ResolveExpressionActivity(usage js.Usage) error {
	errCtx := m.Context()
	err := errCtx.NewError("Internal Error: stack trace should give error before being resolving activity")
	return err
}
