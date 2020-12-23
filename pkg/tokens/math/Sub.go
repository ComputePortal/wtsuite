package math

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	minusWidth      = plusMinusWidth
	extraSubSpacing = extraAddSubSpacing // same on left and right
	minusYOffset    = plusMinusFracYOffset
	minusXOffset    = plusMinusXOffset
	minusThickness  = lineThickness
)

type Sub struct {
	BinaryOp
}

func NewSubOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &Sub{BinaryOp{"-", a, b, newTokenData(ctx)}}, nil
}

func (t *Sub) genMinusPath(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	var b strings.Builder

	d := minusThickness
	l := minusWidth

	x0 := x + minusXOffset
	y0 := y - minusYOffset

	b.WriteString(fmt.Sprintf("M%g %gv%g", x0, y0, d))
	b.WriteString(fmt.Sprintf("h%gv%gh%gz", l, -d, -l))

	if err := scope.BuildMathPath(b.String(), t.Context()); err != nil {
		return nil, err
	}

	bb := boundingbox.NewBB(x0, y0, x0+l, y0+d)

	return bb, nil
}

func (t *Sub) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	bbop, err := t.genMinusPath(scope, bba.Right()+extraSubSpacing, y)
	if err != nil {
		return nil, err
	}

	bbb, err := t.b.GenerateTags(scope, bbop.Right()+extraSubSpacing, y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bba, bbop, bbb), nil
}
