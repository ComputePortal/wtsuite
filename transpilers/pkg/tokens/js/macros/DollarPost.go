package macros

import (
	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type DollarPost struct {
	PostMacro
}

func NewDollarPost(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &DollarPost{newPostMacro(args, ctx)}, nil
}

func (m *DollarPost) Dump(indent string) string {
	return indent + "DollarPost(...)\n"
}

func (m *DollarPost) WriteExpression() string {
	return m.PostMacro.writeExpression(xmlPostHeader.Name())
}

func (m *DollarPost) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if len(args) != 3 {
		return nil, ctx.NewError("Error: expected 3 arguments")
	}

	if !args[0].IsInstanceOf(prototypes.String) {
		return nil, ctx.NewError("Error: expected String for argument 1, got " + args[0].TypeName())
	}

	return m.PostMacro.evalExpression(stack, args[1], args[2])
}

func (m *DollarPost) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(xmlPostHeader, m.Context())

	return m.PostMacro.ResolveExpressionActivity(usage)
}

func (m *DollarPost) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(xmlPostHeader, ns); err != nil {
		return err
	}

	return m.PostMacro.UniqueExpressionNames(ns)
}
