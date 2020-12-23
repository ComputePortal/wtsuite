package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Any struct {
	ValueData
}

func NewAny(ctx context.Context) Value {
	return &Any{ValueData{ctx}}
}

func NewAll(ctx context.Context) Value {
  return &Any{ValueData{ctx}}
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

func IsAny(v_ Value) bool {
  v_ = UnpackContextValue(v_)
  
  if _, ok := v_.(*Any); ok {
    return true
  } else {
    return false
  }
}
