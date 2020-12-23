package macros

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type SharedWorkerPost struct {
	PostMacro
}

func NewSharedWorkerPost(args []js.Expression,
	ctx context.Context) (js.Expression, error) {
	if js.TARGET != "browser" {
		return nil, ctx.NewError("Error: only available if target is browser, (now it is " + js.TARGET + ")")
	}

	if len(args) != 3 {
		return nil, ctx.NewError("Error: expected 3 arguments")
	}

  interfExpr, err := getTypeExpression(args[2])
  if err != nil {
    return nil, err
  }

  return &SharedWorkerPost{newPostMacro(args[0:2], interfExpr, ctx)}, nil
}

func (m *SharedWorkerPost) Dump(indent string) string {
	return indent + "SharedWorkerPost(...)"
}

func (m *SharedWorkerPost) WriteExpression() string {
	return m.PostMacro.writeExpression(sharedWorkerPostHeader.Name())
}

func (m *SharedWorkerPost) EvalExpression() (values.Value, error) {
	ctx := m.Context()

  args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsSharedWorker(args[0]) {
		return nil, ctx.NewError("Error: expected SharedWorker for argument 1, got " + args[0].TypeName())
	}

	return m.PostMacro.evalExpression(args[1])
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
