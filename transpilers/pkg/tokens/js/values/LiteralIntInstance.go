package values

import (
  "../../context"
)

type LiteralIntInstance struct {
  value int

  Instance
}

func NewLiteralIntInstance(interf Interface, i int, ctx context.Context) Value {
  return &LiteralIntInstance{i, newInstance(interf, ctx)}
}

func (v *LiteralIntInstance) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAll(other_) {
    return nil
  } else if other, ok := other_.(*LiteralIntInstance); ok {
    if v.value == other.value {
      return nil
    } else {
      return ctx.NewError("Error: expected other literal int")
    }
  } else {
    return ctx.NewError("Error: not a literal int instance")
  }
}

func (v *LiteralIntInstance) LiteralIntValue() (int, bool) {
  return v.value, true
}
