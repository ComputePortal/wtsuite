package macros

import (
	"fmt"
	"math"
	"strings"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
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
