package math

import (
	"./boundingbox"

	"../context"
)

type BinSymbolOp struct {
	extraSpacingLeft  float64
	extraSpacingRight float64

	symbol *Symbol
	BinaryOp
}

func newBinSymbolOp(left, right float64, symbol *Symbol, a, b Token, ctx context.Context) BinSymbolOp {
	return BinSymbolOp{left, right, symbol, BinaryOp{symbol.symbol, a, b, newTokenData(ctx)}}
}

func NewBinSymbolOp(left, right float64, symbol *Symbol, a, b Token, ctx context.Context) (Token, error) {
	token := newBinSymbolOp(left, right, symbol, a, b, ctx)
	return &token, nil
}

func (t *BinSymbolOp) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bba, err := t.a.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	bbop, err := t.symbol.GenerateTags(scope, bba.Right()+t.extraSpacingLeft, y)
	if err != nil {
		return nil, err
	}

	bbb, err := t.b.GenerateTags(scope, bbop.Right()+t.extraSpacingRight, y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bba, bbop, bbb), nil
}
