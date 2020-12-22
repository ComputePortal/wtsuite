package macros

import (
  "fmt"
  "strings"

  "../values"

  "../../context"
  "../../js"
)

type CastCall struct {
  typeExpr *js.TypeExpression
  Macro
}

func NewCastCall(args []js.Expression, ctx context.Context) (js.Expression, error) {
  // no need to infer the typexpression
  if len(args) != 2 {
    errCtx := ctx
    return nil, errCtx.NewError(fmt.Sprintf("Error: expected 2 arguments, got %d", len(args)))
  }

  typeExpr, err := getTypeExpression(args[1])
  if err != nil {
    return nil, err
  }

  return &CastCall{typeExpr, newMacro(args[0:1], ctx)}, nil
}

func (m *CastCall) Dump(indent string) string {
  var b strings.Builder
  
  b.WriteString(indent)
  b.WriteString("cast(...)")
  b.WriteString("\n")
  b.WriteString(m.args[0].Dump(indent + "  "))
  b.WriteString(m.typeExpr.WriteExpression())
  b.WriteString("\n")

  return b.String()
}

func (m *CastCall) WriteExpression() string {
  return m.args[0].WriteExpression()
}

func (m *CastCall) ResolveExpressionNames(scope js.Scope) error { 
  // last argument
  if err := m.Macro.ResolveExpressionNames(scope); err != nil {
    return err
  }

  if err := m.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  return nil
}

func (m *CastCall) EvalExpression() (values.Value, error) {
  // value of args doesnt matter
  if _, err := m.evalArgs(); err != nil {
    return nil, err
  }

  return m.typeExpr.EvalExpression()
}
