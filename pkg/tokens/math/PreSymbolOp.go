package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type PreSymbolOp struct {
	extraSpacingLeft  float64
	extraSpacingRight float64

	symbol *Symbol
	UnaryOp
}

func newPreSymbolOp(left, right float64, symbol *Symbol, a Token, ctx context.Context) PreSymbolOp {
	return PreSymbolOp{left, right, symbol, UnaryOp{symbol.symbol, a, newTokenData(ctx)}}
}

func (t *PreSymbolOp) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bbSymbol, err := t.symbol.GenerateTags(scope, x+t.extraSpacingLeft, y)
	if err != nil {
		return nil, err
	}

	bba, err := t.a.GenerateTags(scope, bbSymbol.Right()+t.extraSpacingRight, y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bbSymbol, bba), nil
}
