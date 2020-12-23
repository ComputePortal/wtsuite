package values

import (
  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LiteralBooleanInstance struct {
  value bool

  Instance
}

func NewLiteralBooleanInstance(interf Interface, b bool, ctx context.Context) Value {
  return &LiteralBooleanInstance{b, newInstance(interf, ctx)}
}

func (v *LiteralBooleanInstance) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAny(other_) {
    return nil
  } else if other, ok := other_.(*LiteralBooleanInstance); ok {
    if v.value == other.value {
      return nil
    } else {
      return ctx.NewError("Error: have other literal bool")
    }
  } else {
    return ctx.NewError("Error: have " + other_.TypeName() + ", want " + v.TypeName())
  }
}

func (v *LiteralBooleanInstance) LiteralBooleanValue() (bool, bool) {
  return v.value, true
}
