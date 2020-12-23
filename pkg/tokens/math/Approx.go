package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	extraApproxLeftSpacing  = extraEqualsSpacing
	extraApproxRightSpacing = extraEqualsSpacing
)

type Approx struct {
	BinSymbolOp
}

func NewApproxOp(a Token, b Token, ctx context.Context) (Token, error) {
	return &Approx{newBinSymbolOp(extraApproxLeftSpacing, extraApproxRightSpacing, newUnicodeSymbol("~=", 0x2248, ctx), a, b, ctx)}, nil
}
