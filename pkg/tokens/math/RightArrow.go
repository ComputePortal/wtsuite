package math

import (
	"../context"
)

const (
	extraRightArrow1LeftSpacing  = genericBinSymbolSpacing
	extraRightArrow1RightSpacing = genericBinSymbolSpacing
	extraRightArrow2LeftSpacing  = 2 * genericBinSymbolSpacing
	extraRightArrow2RightSpacing = 2 * genericBinSymbolSpacing
)

type RightArrow1 struct {
	BinSymbolOp
}

type RightArrow2 struct {
	BinSymbolOp
}

func NewRightArrow1(a Token, b Token, ctx context.Context) (Token, error) {
	return &RightArrow1{newBinSymbolOp(extraRightArrow1LeftSpacing, extraRightArrow1RightSpacing,
		newUnicodeSymbol("->", 0x2192, ctx), a, b, ctx)}, nil
}

func NewRightArrow2(a Token, b Token, ctx context.Context) (Token, error) {
	return &RightArrow2{newBinSymbolOp(extraRightArrow2LeftSpacing, extraRightArrow2RightSpacing,
		newUnicodeSymbol("=>", 0x21d2, ctx), a, b, ctx)}, nil
}
