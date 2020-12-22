package math

import (
	"../context"
)

const (
	extraGTLeftSpacing  = genericBinSymbolSpacing
	extraGTRightSpacing = genericBinSymbolSpacing
)

type GT struct {
	BinSymbolOp
}

func NewGTOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &GT{newBinSymbolOp(extraGTLeftSpacing, extraGTRightSpacing, newSymbol(">", ctx), a, b, ctx)}, nil
}
