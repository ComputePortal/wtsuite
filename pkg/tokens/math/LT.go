package math

import (
	"../context"
)

const (
	extraLTLeftSpacing  = genericBinSymbolSpacing
	extraLTRightSpacing = genericBinSymbolSpacing
)

type LT struct {
	BinSymbolOp
}

func NewLTOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &LT{newBinSymbolOp(extraLTLeftSpacing, extraLTRightSpacing, newSymbol("<", ctx), a, b, ctx)}, nil
}
