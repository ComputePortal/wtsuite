package macros

import (
	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type SharedWorkerPost struct {
	PostMacro
}

func NewSharedWorkerPost(args []js.Expression,
	ctx context.Context) (js.Expression, error) {
	if js.TARGET != "browser" {
		return nil, ctx.NewError("Error: only available if target is browser, (now it is " + js.TARGET + ")")
	}

	return &SharedWorkerPost{newPostMacro(args, ctx)}, nil
}

func (m *SharedWorkerPost) Dump(indent string) string {
	return indent + "SharedWorkerPost(...)"
}

func (m *SharedWorkerPost) WriteExpression() string {
	return m.PostMacro.writeExpression(sharedWorkerPostHeader.Name())
}

func (m *SharedWorkerPost) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if len(args) != 3 {
		return nil, ctx.NewError("Error: expected 3 arguments")
	}

	if !args[0].IsInstanceOf(prototypes.SharedWorker) {
		return nil, ctx.NewError("Error: expected SharedWorker for argument 1, got " + args[0].TypeName())
	}

	return m.PostMacro.evalExpression(stack, args[1], args[2])
}

func (m *SharedWorkerPost) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(sharedWorkerPostHeader, m.Context())

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *SharedWorkerPost) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(sharedWorkerPostHeader, ns); err != nil {
		return err
	}

	return m.Macro.UniqueExpressionNames(ns)
}
