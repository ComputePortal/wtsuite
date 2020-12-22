package macros

import (
  "strings"

  "../prototypes"
  "../values"

  "../../context"
  "../../js"
)

type MathAdvanceWidth struct {
	Macro
}

func NewMathAdvanceWidth(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return &MathAdvanceWidth{newMacro(args, ctx)}, nil
}

func (m *MathAdvanceWidth) Dump(indent string) string {
	return indent + "MathAdvanceWidth(...)"
}

func (m *MathAdvanceWidth) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

  if !prototypes.IsInt(args[0]) {
		return nil, ctx.NewError("Error: expected Int, got " + args[0].TypeName())
  }

  return prototypes.NewInt(ctx), nil
}

func (m *MathAdvanceWidth) WriteExpression() string {
  var b strings.Builder

  b.WriteString(mathFontHeader.Name())
  b.WriteString(".advanceWidth(")
  b.WriteString(m.args[0].WriteExpression())
  b.WriteString(")")

  return b.String()
}

func (m *MathAdvanceWidth) ResolveExpressionActivity(usage js.Usage) error {
  ResolveHeaderActivity(mathFontHeader, m.Context())

  return m.Macro.ResolveExpressionActivity(usage)
}
