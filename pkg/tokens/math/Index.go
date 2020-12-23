package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	indexLevelDecrease   = powLevelDecrease
	extraIndexSpacing    = extraPowSpacing
	indexSubscriptOffset = -0.2
)

type Index struct {
	BinaryOp
}

func NewIndexOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &Index{BinaryOp{"_", a, b, newTokenData(ctx)}}, nil
}

func (t *Index) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	dx := bba.Right() + extraIndexSpacing
	dy := y - indexSubscriptOffset

	subScope := scope.NewSubScope()
	subScope.Transform(dx, dy, indexLevelDecrease, indexLevelDecrease)
	bbb, err := t.b.GenerateTags(subScope, 0.0, 0.0)
	if err != nil {
		return nil, err
	}

	bbb = bbb.Scale(indexLevelDecrease, indexLevelDecrease)
	bbb = bbb.Translate(dx, dy)

	return boundingbox.Merge(bba, bbb), nil
}
