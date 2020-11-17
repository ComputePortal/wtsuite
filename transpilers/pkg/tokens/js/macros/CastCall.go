package macros

import (
  "fmt"
  "strings"

  "../values"

  "../../context"
  "../../js"
)

type CastCall struct {
  Macro
}

func NewCastCall(args []js.Expression, ctx context.Context) (js.Expression, error) {
  // no need to infer the typexpression
  if len(args) != 2 {
    errCtx := ctx
    return nil, errCtx.NewError(fmt.Sprintf("Error: expected 2 arguments, got %d", len(args)))
  }

  return &CastCall{newMacro(args, ctx)}, nil
}

func (m *CastCall) Dump(indent string) string {
  var b strings.Builder
  
  b.WriteString(indent)
  b.WriteString("cast(...)")
  b.WriteString("\n")
  b.WriteString(m.args[0].Dump(indent + "  "))
  b.WriteString(m.args[1].Dump(indent + "  "))
  b.WriteString("\n")

  return b.String()
}

func (m *CastCall) WriteExpression() string {
  return m.args[0].WriteExpression()
}

func (m *CastCall) EvalExpression() (values.Value, error) {
  // value of first argument doesnt matter
  if _, err := m.args[0].EvalExpression(); err != nil {
    return nil, err
  }

  return m.args[1].EvalExpression()
}
