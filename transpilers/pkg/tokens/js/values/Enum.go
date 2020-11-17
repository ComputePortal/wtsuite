package values

import (
	"../../context"
)

type Enum struct {
  proto Prototype
	ValueData
}

func NewEnum(proto Prototype, ctx context.Context) *Enum {
	return &Enum{proto, ValueData{ctx}}
}

func (v *Enum) TypeName() string {
  return v.proto.Name()
}

func (v *Enum) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAll(other_) {
    return nil
  } else if other, ok := other_.(*Enum); ok && other == v {
    if other != v {
      return ctx.NewError("Error: different enum")
    }

    return nil
  } else {
    return ctx.NewError("Error: not an enum")
  }
}

func (v *Enum) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: enum can't be constructed")
}

func (v *Enum) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't call a enum")
}

func (v *Enum) GetMember(key string, includePrivate bool,
  ctx context.Context) (Value, error) {
  return v.proto.GetClassMember(key, includePrivate, ctx)
}

func (v *Enum) SetMember(key string, includePrivate bool, arg Value,
  ctx context.Context) error {
  return ctx.NewError("Error: can't set static enum members")
}
