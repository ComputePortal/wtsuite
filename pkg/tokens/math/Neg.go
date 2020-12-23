package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	extraNegLeftSpacing  = 0.0
	extraNegRightSpacing = genericBinSymbolSpacing
)

type Neg struct {
	PreSymbolOp
}

// XXX: should we create minus sign like Sub?
func NewNegOp(a Token, ctx context.Context) (Token, error) {
	return &Neg{newPreSymbolOp(extraNegLeftSpacing, extraNegRightSpacing, newSymbol("-", ctx), a, ctx)}, nil
}
