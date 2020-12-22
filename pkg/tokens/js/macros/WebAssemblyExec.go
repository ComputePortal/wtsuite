package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type WebAssemblyExec struct {
	Macro
}

func NewWebAssemblyExec(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	return &WebAssemblyExec{newMacro(args, ctx)}, nil
}

func (m *WebAssemblyExec) Dump(indent string) string {
	return indent + "WebAssemblyExec(url, env)"
}

func (m *WebAssemblyExec) WriteExpression() string {
	var b strings.Builder

	b.WriteString("(()=>{return new Promise((o,n)=>{")
	b.WriteString("var e=")
	b.WriteString(m.args[1].WriteExpression())
	b.WriteString(";WebAssembly.instantiateStreaming(fetch(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString("),{env:e")
	b.WriteString("}).then((r)=>{e.heapOffset=r.instance.exports.__heap_base;r.instance.exports.main();o()}).catch((e)=>{n(e)})})})()")

	return b.String()
}

func (m *WebAssemblyExec) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsString(args[0]) {
		errCtx := m.args[0].Context()
		return nil, errCtx.NewError("Error: expected string, got " + args[0].TypeName())
	}

	if !prototypes.IsWebAssemblyEnv(args[1]) {
		errCtx := m.args[1].Context()
		return nil, errCtx.NewError("Error: expected WebAssemblyEnv, got " + args[1].TypeName() + " (hint: wrap WebAssemblyFS by WebAssemblyEnv)")
	}

	return prototypes.NewVoidPromise(ctx), nil
}
