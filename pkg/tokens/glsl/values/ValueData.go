package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ValueData struct {
  ctx context.Context
}

func newValueData(ctx context.Context) ValueData {
  return ValueData{ctx}
}

func (v *ValueData) Context() context.Context {
  return v.ctx
}

func (v *ValueData) Instantiate(ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't instantiate")
}

func (v *ValueData) LiteralIntValue() (int, bool) {
  return 0, false
}
