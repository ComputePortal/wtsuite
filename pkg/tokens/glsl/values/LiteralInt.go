package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LiteralInt struct {
  i int
  Scalar
}

func NewLiteralInt(i int, ctx context.Context) *LiteralInt {
  return &LiteralInt{i, newScalar("int", ctx)}
}

func (v *LiteralInt) LiteralIntValue() (int, bool) {
  return v.i, true
}

func AssertLiteralInt(v_ Value) (*LiteralInt, error) {
  errCtx := v_.Context()

  v_ = UnpackContextValue(v_)

  if v, ok := v_.(*LiteralInt); ok {
    return v, nil
  } else {
    return nil, errCtx.NewError("Error: not a literal int")
  }
}
