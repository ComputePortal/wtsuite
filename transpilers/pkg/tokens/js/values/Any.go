package values

import (
	"../../context"
)

type Any struct {
  isAll bool 

	ValueData
}

func NewAny(ctx context.Context) Value {
	return &Any{false, ValueData{ctx}}
}

func NewAll(ctx context.Context) Value {
  return &Any{true, ValueData{ctx}}
}

func (v *Any) TypeName() string {
  return "any"
}

func (v *Any) Check(other Value, ctx context.Context) error {
  return nil
}

func (v *Any) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
	return NewAny(ctx), nil
}

func (v *Any) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  if preferMethod {
    return nil, nil
  } else {
    return NewAny(ctx), nil
  }
}

func (v *Any) GetMember(key string, includePrivate bool, ctx context.Context) (Value, error) {
  return NewAny(ctx), nil
}

func (v *Any) SetMember(key string, includePrivate bool, arg Value, ctx context.Context) error {
  return nil
}

func IsAll(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  v, ok := v_.(*Any)
  if ok {
    return v.isAll
  } else {
    return false
  }
}
