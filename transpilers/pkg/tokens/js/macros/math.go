package macros

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
	"../../math/serif"
)

type Convert struct {
	// result = c0 + c1*input
	c0 float64 // constant
	c1 float64 // scale factor
	Macro
}

func newConvert(c0 float64, c1 float64, args []js.Expression,
	ctx context.Context) Convert {
	return Convert{c0, c1, newMacro(args, ctx)}
}

func (m *Convert) WriteExpression() string {
	var b strings.Builder

	b.WriteString("(")

	if m.c0 != 0.0 {
		b.WriteString(fmt.Sprintf("%.08f+", m.c0))
	}

	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	if m.c1 != 1.0 {
		b.WriteString(fmt.Sprintf("*%.08f", m.c1))
	}

	b.WriteString(")")

	return b.String()
}

func (m *Convert) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if !prototypes.IsNumber(args[0]) {
		return nil, ctx.NewError("Error: expected Number argument, got " + args[0].TypeName())
	}

	return prototypes.NewNumber(ctx), nil
}

type DegToRad struct {
	Convert
}

func NewDegToRad(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &DegToRad{newConvert(0.0, math.Pi/180.0, args, ctx)}, nil
}

func (m *DegToRad) Dump(indent string) string {
	return indent + "DegToRad(...)"
}

type RadToDeg struct {
	Convert
}

func NewRadToDeg(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &RadToDeg{newConvert(0.0, 180.0/math.Pi, args, ctx)}, nil
}

func (m *RadToDeg) Dump(indent string) string {
	return indent + "RadToDeg(...)"
}

type MathAdvanceWidth struct {
	aw int
	Macro
}

func NewMathAdvanceWidth(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &MathAdvanceWidth{0, newMacro(args, ctx)}, nil
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

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if litInt, ok := args[0].LiteralIntValue(); ok {
    aw, okInner := serif.AdvanceWidths[litInt]
    if !okInner {
      err := ctx.NewError(fmt.Sprintf("Error: advance width for %d not found", litInt))
      panic(err)
    }

    m.aw = aw

		return prototypes.NewLiteralInt(m.aw, ctx), nil
	} else {
		return nil, ctx.NewError("Error: expected a literal int, got " + args[0].TypeName())
	}
}

func (m *MathAdvanceWidth) WriteExpression() string {
	return strconv.Itoa(m.aw)
}

type MathBoundingBox struct {
	bb []float64
	Macro
}

func NewMathBoundingBox(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &MathBoundingBox{make([]float64, 4), newMacro(args, ctx)}, nil
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

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if litInt, ok := args[0].LiteralIntValue(); ok {
		bb, okInner := serif.Bounds[litInt]

    if !okInner {
      err := ctx.NewError(fmt.Sprintf("Error: bounds for char %d not found", litInt))
      panic(err)
    }

		m.bb[0] = bb.Left()
		m.bb[1] = bb.Right()
		m.bb[2] = bb.Top()
		m.bb[3] = bb.Bottom()

		return prototypes.NewArray(prototypes.NewNumber(ctx), ctx), nil
	} else {
		return nil, ctx.NewError("Error: expected a literal int, got " + args[0].TypeName())
	}
}

func (m *MathBoundingBox) WriteExpression() string {
	return fmt.Sprintf("[%g,%g,%g,%g]", m.bb[0], m.bb[1], m.bb[2], m.bb[3])
}

type MathFormatMetrics struct {
	Macro
}

func NewMathFormatMetrics(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &MathFormatMetrics{newMacro(args, ctx)}, nil
}

func (m *MathFormatMetrics) Dump(indent string) string {
	return indent + "MathFormatMetrics(...)"
}

func (m *MathFormatMetrics) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected no arguments")
	}

	// TODO: add literal value entries
	return prototypes.NewObject(make(map[string]values.Value), ctx), nil
}

func (m *MathFormatMetrics) WriteExpression() string {
	// TODO: add literal value entries
	return "{}"
}
