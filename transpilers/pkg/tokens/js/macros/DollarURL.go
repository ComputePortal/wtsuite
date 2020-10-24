package macros

import (
	"../"

	"../prototypes"
	"../values"

	"../../context"

	"../../../files"
)

type DollarURL struct {
	url string
	BrowserMacro
}

func NewDollarURL(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &DollarURL{"", newBrowserMacro(args, ctx)}, nil
}

func (m *DollarURL) Dump(indent string) string {
	return indent + "DollarURL(...)"
}

func (m *DollarURL) WriteExpression() string {
	return "'" + m.url + "'"
}

func (m *DollarURL) ResolveExpressionNames(scope js.Scope) error {
	ctx := m.Context()

	if len(m.args) != 1 {
		return ctx.NewError("Error: expeced 1 argument")
	}

	if _, ok := m.args[0].(*js.LiteralString); !ok {
		return ctx.NewError("Error: expected a literal string")
	}

	return m.BrowserMacro.ResolveExpressionNames(scope)
}

func (m *DollarURL) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	litStr, _ := m.args[0].(*js.LiteralString)
	absPath, err := files.Search(ctx.Path(), litStr.Value())
	if err != nil {
		errCtx := litStr.Context()
		return nil, errCtx.NewError("Error: file " + litStr.Value() + " not found")
	}

	vif := stack.GetViewInterface(absPath)
	if vif == nil {
		return nil, ctx.NewError("Error: " + litStr.Value() + " not a view")
	}

	url := vif.GetURL()

	if m.url != "" && m.url != url {
		panic("internal error")
	}

	m.url = url

	return prototypes.NewLiteralString(url, m.Context()), nil
}
