package math

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	plusWidth       = plusMinusWidth // same as minus
	extraAddSpacing = extraAddSubSpacing
	plusYOffset     = plusMinusFracYOffset
	plusXOffset     = plusMinusXOffset
	plusThickness   = lineThickness
)

type Add struct {
	BinaryOp
}

func NewAddOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &Add{BinaryOp{"+", a, b, newTokenData(ctx)}}, nil
}

func (t *Add) genPlusPath(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	var b strings.Builder

	d := plusThickness
	l := 0.5 * (plusWidth - plusThickness)

	x0 := x + plusXOffset
	y0 := y - plusYOffset

	b.WriteString(fmt.Sprintf("M%g %gv%g", x0, y0, d))
	b.WriteString(fmt.Sprintf("h%gv%gh%gv%gh%g", l, l, d, -l, l))
	b.WriteString(fmt.Sprintf("v%gh%gv%gh%gv%gz", -d, -l, -l, -d, l))

	if err := scope.BuildMathPath(b.String(), t.Context()); err != nil {
		return nil, err
	}

	return boundingbox.NewBB(x0, y0-l, x0+plusWidth, y0+d+l), nil
}

func (t *Add) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	bbop, err := t.genPlusPath(scope, bba.Right()+extraAddSpacing, y)
	if err != nil {
		return nil, err
	}

	bbb, err := t.b.GenerateTags(scope, bbop.Right()+extraAddSpacing, y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bba, bbop, bbb), nil
}
