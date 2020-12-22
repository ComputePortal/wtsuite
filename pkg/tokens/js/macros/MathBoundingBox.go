package macros

import (
  "strings"

  "../prototypes"
  "../values"

  "../../context"
  "../../js"
)

type MathBoundingBox struct {
	Macro
}

func NewMathBoundingBox(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return &MathBoundingBox{newMacro(args, ctx)}, nil
}

func (m *MathBoundingBox) Dump(indent string) string {
	return indent + "MathBoundingBox(...)"
}

func (m *MathBoundingBox) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

  if !prototypes.IsInt(args[0]) {
		return nil, ctx.NewError("Error: expected Int, got " + args[0].TypeName())
  }

  return prototypes.NewArray(prototypes.NewNumber(ctx), ctx), nil
}

func (m *MathBoundingBox) WriteExpression() string {
  var b strings.Builder

  b.WriteString(mathFontHeader.Name())
  b.WriteString(".boundingBox(")
  b.WriteString(m.args[0].WriteExpression())
  b.WriteString(")")

  return b.String()
}

func (m *MathBoundingBox) ResolveExpressionActivity(usage js.Usage) error {
  ResolveHeaderActivity(mathFontHeader, m.Context())

  return m.Macro.ResolveExpressionActivity(usage)
}
