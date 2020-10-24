package math

import (
	"../context"
)

const (
	extraNELeftSpacing  = extraEqualsSpacing
	extraNERightSpacing = extraEqualsSpacing
)

type NE struct {
	BinSymbolOp
}

func NewNEOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &NE{newBinSymbolOp(extraNELeftSpacing, extraNERightSpacing, newUnicodeSymbol("!=", 0x2260, ctx), a, b, ctx)}, nil
}
