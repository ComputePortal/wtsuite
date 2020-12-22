package math

import (
	"fmt"
	"strings"

	"./boundingbox"

	"../context"
)

const (
	equalsWidth         = plusMinusWidth
	extraEqualsSpacing  = extraAddSubSpacing // same on left and right
	equalsYTopOffset    = 0.386              // subtract from y!
	equalsYBottomOffset = 0.186
	equalsXOffset       = plusMinusXOffset
	equalsThickness     = lineThickness
)

type Equals struct {
	BinaryOp
}

func NewEqualsOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &Equals{BinaryOp{"=", a, b, newTokenData(ctx)}}, nil
}

func (t *Equals) genEqualsPath(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	var b strings.Builder

	d := equalsThickness
	l := equalsWidth

	x0 := x + equalsXOffset
	y0 := y - equalsYTopOffset

	x1 := x0
	y1 := y - equalsYBottomOffset

	b.WriteString(fmt.Sprintf("M%g %gv%g", x0, y0, d))
	b.WriteString(fmt.Sprintf("h%gv%gh%gz", l, -d, -l))

	b.WriteString(fmt.Sprintf("M%g %gv%g", x1, y1, d))
	b.WriteString(fmt.Sprintf("h%gv%gh%gz", l, -d, -l))

	if err := scope.BuildMathPath(b.String(), t.Context()); err != nil {
		return nil, err
	}

	bb := boundingbox.NewBB(x0, y0, x1+l, y1+d)

	return bb, nil
}

func (t *Equals) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	bbop, err := t.genEqualsPath(scope, bba.Right()+extraEqualsSpacing, y)
	if err != nil {
		return nil, err
	}

	bbb, err := t.b.GenerateTags(scope, bbop.Right()+extraEqualsSpacing, y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bba, bbop, bbb), nil
}
