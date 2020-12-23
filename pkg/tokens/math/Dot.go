package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	dotScale      = 0.5
	dotVerSpacing = -0.1 // relative to top of center of bb
	dotHorOffset  = 0.15
)

type Dot struct {
	BinaryOp // abuse b for the dot symbol
}

func NewDot(a Token, ctx context.Context) (Token, error) {
	b := newUnicodeSymbol(".", 0x2219, ctx)

	return &Dot{BinaryOp{".", a, b, newTokenData(ctx)}}, nil
}

func (t *Dot) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	subScope := scope.NewSubScope()

	bbDot, err := t.b.GenerateTags(subScope, dotHorOffset, 0.0)
	if err != nil {
		return nil, err
	}

	dx := bba.Left() + 0.5*bba.Width() - 0.5*bbDot.Width()
	dy := bba.Top() + dotVerSpacing

	subScope.Transform(dx, dy, dotScale, dotScale)

	bbDot = bbDot.Scale(dotScale, dotScale)
	bbDot = bbDot.Translate(dx, dy)

	return boundingbox.Merge(bba, bbDot), nil
}
