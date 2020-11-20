package macros

import (
	"../values"

  "../../context"
  "../../js"
)

type ToInstance struct {
  interfExpr *js.TypeExpression
	interf     values.Interface // collected during resolve names stage
  Macro
}

func newToInstance(args []js.Expression, interfExpr *js.TypeExpression, ctx context.Context) ToInstance {
	return ToInstance{interfExpr, nil, newMacro(args, ctx)}
}

func getTypeExpression(expr_ js.Expression) (*js.TypeExpression, error) {
  expr, err := js.GetTypeExpression(expr_)
  if err != nil {
    return nil, err
  }

  if expr == nil {
    errCtx := expr_.Context()
    return nil, errCtx.NewError("Error: could never be an interface")
  }

  return expr, nil
}

func (m *ToInstance) ResolveExpressionNames(scope js.Scope) error {
  // last argument
  if err := m.Macro.ResolveExpressionNames(scope); err != nil {
    return err
  }

  if err := m.interfExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  m.interf = m.interfExpr.GetInterface() 
  if m.interf == nil {
    errCtx := m.interfExpr.Context()
    return errCtx.NewError("Error: not an interface or a prototype")
  }

  return nil
}

func (m *ToInstance) ResolveExpressionActivity(usage js.Usage) error {
  if err := m.Macro.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  if !m.interf.IsUniversal() {
    if _, ok := m.interf.(values.Prototype); ok {
      errCtx := m.Context()
      return errCtx.NewError("Error: prototype " + m.interf.Name() + " is not universal (hint: use 'universe')")
    } else {
      errCtx := m.Context()
      return errCtx.NewError("Error: interface " + m.interf.Name() + " is not universal (hint: use 'universe' for all implementations)")
    }
  }

  return nil
}
