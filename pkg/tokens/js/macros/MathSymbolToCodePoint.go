package macros

import (
  "strings"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
  "github.com/computeportal/wtsuite/pkg/tokens/js"
  "github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type MathSymbolToCodePoint struct {
	Macro
}

func NewMathSymbolToCodePoint(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return &MathSymbolToCodePoint{newMacro(args, ctx)}, nil
}

func (m *MathSymbolToCodePoint) Dump(indent string) string {
	return indent + "MathSymbolToCodePoint(...)"
}

func (m *MathSymbolToCodePoint) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

  if !prototypes.IsString(args[0]) {
		return nil, ctx.NewError("Error: expected String, got " + args[0].TypeName())
  }

  return prototypes.NewInt(ctx), nil
}

func (m *MathSymbolToCodePoint) WriteExpression() string {
  var b strings.Builder

  b.WriteString(mathFontHeader.Name())
  b.WriteString(".symbolToCodePoint(")
  b.WriteString(m.args[0].WriteExpression())
  b.WriteString(")")

  return b.String()
}

func (m *MathSymbolToCodePoint) ResolveExpressionActivity(usage js.Usage) error {
  ResolveHeaderActivity(mathFontHeader, m.Context())

  return m.Macro.ResolveExpressionActivity(usage)
}
