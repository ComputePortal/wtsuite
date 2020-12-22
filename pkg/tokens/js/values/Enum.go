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

  if IsAny(other_) {
    return nil
  } else if other, ok := other_.(*Enum); ok && other == v {
    return nil
  } 

  return ctx.NewError("Error: have " + other_.TypeName() + ", want " + v.TypeName())
}

func (v *Enum) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a constructor")
}

func (v *Enum) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *Enum) GetMember(key string, includePrivate bool,
  ctx context.Context) (Value, error) {
  if res, err := v.proto.GetClassMember(key, includePrivate, ctx); err != nil {
    return nil, err
  } else if res != nil {
    return res, nil
  } else {
    return nil, ctx.NewError("Error: " + v.proto.Name() + "." + key + " not found")
  }
}

func (v *Enum) SetMember(key string, includePrivate bool, arg Value,
  ctx context.Context) error {
  return ctx.NewError("Error: can't set static enum members")
}
