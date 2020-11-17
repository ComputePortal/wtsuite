package macros

import (
	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type XMLHttpRequestPost struct {
	PostMacro
}

func NewXMLHttpRequestPost(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &XMLHttpRequestPost{newPostMacro(args, ctx)}, nil
}

func (m *XMLHttpRequestPost) Dump(indent string) string {
	return indent + "XMLHttpRequestPost(...)\n"
}

func (m *XMLHttpRequestPost) WriteExpression() string {
	return m.PostMacro.writeExpression(xmlPostHeader.Name())
}

func (m *XMLHttpRequestPost) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if len(args) != 3 {
		return nil, ctx.NewError("Error: expected 3 arguments")
	}

	if !prototypes.IsString(args[0]) {
		return nil, ctx.NewError("Error: expected String for argument 1, got " + args[0].TypeName())
	}

	return m.PostMacro.evalExpression(args[1], args[2])
}

func (m *XMLHttpRequestPost) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(xmlPostHeader, m.Context())

	return m.PostMacro.ResolveExpressionActivity(usage)
}

func (m *XMLHttpRequestPost) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(xmlPostHeader, ns); err != nil {
		return err
	}

	return m.PostMacro.UniqueExpressionNames(ns)
}
