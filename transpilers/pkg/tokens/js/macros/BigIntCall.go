package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type BigIntCall struct {
	Macro
}

func NewBigIntCall(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &BigIntCall{newMacro(args, ctx)}, nil
}

func (m *BigIntCall) Dump(indent string) string {
	return indent + "BigIntCall(...)"
}

func (m *BigIntCall) WriteExpression() string {
	// XXX: should everything be wrapped in additional parentheses?
	var b strings.Builder

	b.WriteString("BigInt(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (m *BigIntCall) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()
	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if err := prototypes.CheckInputs(&prototypes.Any{}, args, ctx); err != nil {
		return nil, err
	}

	return prototypes.NewBigInt(ctx), nil
}
