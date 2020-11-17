package macros

import (
	"../"

	"../values"

	"../../context"
)

type BrowserMacro struct {
	Macro
}

func newBrowserMacro(args []js.Expression, ctx context.Context) BrowserMacro {
	return BrowserMacro{newMacro(args, ctx)}
}

func (m *BrowserMacro) evalArgs(ctx context.Context) ([]values.Value, error) {
	if js.TARGET != "browser" {
		return nil, ctx.NewError("Error: illegal if TARGET==" + js.TARGET)
	}

	return m.Macro.evalArgs()
}

func (m *BrowserMacro) ResolveExpressionActivity(usage js.Usage) error {
	if js.TARGET != "browser" {
		return nil
		//errCtx := m.Context()
		//err := errCtx.NewError("Internal Error: browser macro shouldnt be resolved for " + js.TARGET)
		//return err
	}

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *BrowserMacro) UniqueExpressionNames(ns js.Namespace) error {
	if js.TARGET != "browser" {
		return nil
		//panic("shouldnt be resolved")
	}

	return m.Macro.UniqueExpressionNames(ns)
}
