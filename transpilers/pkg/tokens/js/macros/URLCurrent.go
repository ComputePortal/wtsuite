package macros

import (
	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type URLCurrent struct {
	Macro
}

func NewURLCurrent(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &URLCurrent{newMacro(args, ctx)}, nil
}

func (m *URLCurrent) Dump(indent string) string {
	return indent + "URLCurrentMacro(...)"
}

func (m *URLCurrent) WriteExpression() string {
	return "(new URL(window.location.href))"
}

func (m *URLCurrent) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	return prototypes.NewURL(ctx), nil
}
