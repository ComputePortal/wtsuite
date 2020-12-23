package macros

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"

	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type URLCurrent struct {
	Macro
}

func NewURLCurrent(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

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

	if _, err := m.evalArgs(); err != nil {
		return nil, err
	}

	return prototypes.NewURL(ctx), nil
}
