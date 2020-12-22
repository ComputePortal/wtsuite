package math

import (
	"./boundingbox"

	"../context"
)

const (
	powLevelDecrease  = 0.8
	extraPowSpacing   = 0.05
	powExponentOffset = 0.5
)

type Pow struct {
	BinaryOp
}

func newPow(a Token, b Token, ctx context.Context) Token {
	return &Pow{BinaryOp{"^", a, b, newTokenData(ctx)}}
}

func NewPowOp(a Token, b Token, ctx context.Context) (Token, error) {
	return newPow(a, b, ctx), nil
}

func (t *Pow) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}
	subScope := scope.NewSubScope()
	bbb, err := t.b.GenerateTags(subScope, 0.0, 0.0)
	if err != nil {
		return nil, err
	}

	dx := bba.Right() + extraPowSpacing
	dy := y - powExponentOffset
	subScope.Transform(dx, dy, powLevelDecrease, powLevelDecrease)

	bbb = bbb.Scale(powLevelDecrease, powLevelDecrease)
	bbb = bbb.Translate(dx, dy)

	return boundingbox.Merge(bba, bbb), nil
}
