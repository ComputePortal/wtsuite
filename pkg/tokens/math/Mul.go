package math

import (
	"./boundingbox"

	"../context"
)

const (
	extraMulSpacing          = 0.10 // used to 0.05
	extraMulSpacingNoSymbols = 0.2
)

type Mul struct {
	BinaryOp
}

func NewMulOp(a Token, b Token, ctx context.Context) (*Mul, error) {
	return &Mul{BinaryOp{"*", a, b, newTokenData(ctx)}}, nil
}

func newMul(a Token, b Token, ctx context.Context) *Mul {
	m, err := NewMulOp(a, b, ctx)
	if err != nil {
		panic(err)
	}

	return m
}

func (t *Mul) Spacing() float64 {
	spacing := extraMulSpacingNoSymbols
	if IsSymbol(t.a) || IsWord(t.a) || IsFloat(t.a) || IsSymbol(t.b) || IsWord(t.b) || IsFloat(t.b) {
		spacing = extraMulSpacing
	}

	return spacing
}

func (t *Mul) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	bbb, err := t.b.GenerateTags(scope, bba.Right()+t.Spacing(), y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bba, bbb), nil
}
