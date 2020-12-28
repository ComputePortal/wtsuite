package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// intended for LiteralInt, LiteralBool, LiteralFloat and LiteralString
type LiteralData struct {
	TokenData
}

func newLiteralData(ctx context.Context) LiteralData {
	return LiteralData{TokenData{ctx}}
}
