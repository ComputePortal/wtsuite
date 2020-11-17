package values

import (
  "../../context"
)

type LiteralStringInstance struct {
  value string

  Instance
}

func NewLiteralStringInstance(interf Interface, str string, ctx context.Context) Value {
  return &LiteralStringInstance{str, newInstance(interf, ctx)}
}

func (v *LiteralStringInstance) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAll(other_) {
    return nil
  } else if other, ok := other_.(*LiteralStringInstance); ok {
    if v.value == other.value {
      return nil
    } else {
      return ctx.NewError("Error: expected literal string \"" + v.value + "\", got \"" + other.value + "\"")
    }
  } else {
    return ctx.NewError("Error: not a literal string instance")
  }
}

func (v *LiteralStringInstance) LiteralStringValue() (string, bool) {
  return v.value, true
}
