package math

import (
	"../context"
)

const (
	extraLELeftSpacing  = genericBinSymbolSpacing
	extraLERightSpacing = genericBinSymbolSpacing
)

type LE struct {
	BinSymbolOp
}

func NewLEOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &LE{newBinSymbolOp(extraLELeftSpacing, extraLERightSpacing, newUnicodeSymbol("<=", 0x2264, ctx), a, b, ctx)}, nil
}
