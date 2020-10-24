package math

import (
	"../context"
)

const (
	extraDivLeftSpacing  = 0.0
	extraDivRightSpacing = 0.0
)

type Div struct {
	BinSymbolOp
}

func NewDivOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &Div{newBinSymbolOp(extraDivLeftSpacing, extraDivRightSpacing, newSymbol("/", ctx), a, b, ctx)}, nil
}
