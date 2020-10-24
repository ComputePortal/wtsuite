package math

import (
	"../context"
)

const (
	extraGELeftSpacing  = genericBinSymbolSpacing
	extraGERightSpacing = genericBinSymbolSpacing
)

type GE struct {
	BinSymbolOp
}

func NewGEOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &GE{newBinSymbolOp(extraGELeftSpacing, extraGERightSpacing, newUnicodeSymbol(">=", 0x2265, ctx), a, b, ctx)}, nil
}
