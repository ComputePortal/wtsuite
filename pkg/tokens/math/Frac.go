package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	extraNumeratorSpacing   = 0.15 // in vertical direction
	extraDenominatorSpacing = 0.15 // in vertical direction
	extraFracSpacing        = 0.0  // in horizontal direction
	nestedFracLevelDecrease = 0.8
	fracThickness           = lineThickness
	fracYOffset             = plusMinusFracYOffset
)

type Frac struct {
	BinaryOp
}

func NewFracOp(a Token, b Token, ctx context.Context) (*Frac, error) {
	return &Frac{BinaryOp{"/", a, b, newTokenData(ctx)}}, nil
}

func newFrac(a Token, b Token, ctx context.Context) *Frac {
	f, err := NewFracOp(a, b, ctx)
	if err != nil {
		panic(err)
	}

	return f
}

func (t *Frac) genFracPath(scope Scope, x float64, y float64, l float64) (boundingbox.BB, error) {
	var b strings.Builder

	d := fracThickness

	x0 := x
	y0 := y - fracYOffset

	b.WriteString(fmt.Sprintf("M%g %gv%g", x0, y0, d))
	b.WriteString(fmt.Sprintf("h%gv%gh%gz", l, -d, -l))

	if err := scope.BuildMathPath(b.String(), t.Context()); err != nil {
		return nil, err
	}

	return boundingbox.NewBB(x0, y0, x0+l, y0+fracThickness), nil
}

func (t *Frac) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	subScopeA := scope.NewSubScope() // transform the subScope later

	bba, err := t.a.GenerateTags(subScopeA, 0.0, 0.0)
	if err != nil {
		return nil, err
	}

	subScopeB := scope.NewSubScope()

	bbb, err := t.b.GenerateTags(subScopeB, 0.0, 0.0)
	if err != nil {
		return nil, err
	}

	l := math.Max(bba.Width(), bbb.Width())
	bbfrac, err := t.genFracPath(scope, x+extraFracSpacing, y, l)
	if err != nil {
		return nil, err
	}

	dxBBA := bbfrac.Left() + bbfrac.Width()*0.5 - bba.Width()*0.5 - bba.Left()
	dyBBA := bbfrac.Top() - extraNumeratorSpacing - bba.Bottom()

	dxBBB := bbfrac.Left() + bbfrac.Width()*0.5 - bbb.Width()*0.5 - bbb.Left()
	dyBBB := bbfrac.Bottom() + extraDenominatorSpacing - bbb.Top()

	subScopeA.Transform(dxBBA, dyBBA, 1.0, 1.0)
	bba = bba.Translate(dxBBA, dyBBA)

	subScopeB.Transform(dxBBB, dyBBB, 1.0, 1.0)
	bbb = bbb.Translate(dxBBB, dyBBB)

	return boundingbox.Merge(bba, bbb, bbfrac), nil
}
