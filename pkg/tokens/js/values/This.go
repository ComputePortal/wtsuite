package values

import (
  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

// so we can assert each property is touched during construction
type This struct {
  this Value
  touched map[string]Value
  ValueData
}

func NewThis(this Value, ctx context.Context) *This {
  return &This{
    this,
    make(map[string]Value),
    ValueData{ctx},
  }
}

func (v *This) TypeName() string {
  return v.this.TypeName()
}

func (v *This) Check(other_ Value, ctx context.Context) error {
  return v.this.Check(other_, ctx)
}

func (v *This) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't construct this")
}

func (v *This) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't call this")
}

func (v *This) GetMember(key string, includePrivate bool, ctx context.Context) (Value, error) {
  return v.this.GetMember(key, includePrivate, ctx)
}

func (v *This) SetMember(key string, includePrivate bool, arg Value, ctx context.Context) error {
  if err := v.this.SetMember(key, includePrivate, arg, ctx); err != nil {
    return err
  }

  v.touched[key] = arg
  return nil
}

func (v *This) AssertTouched(key string, ctx context.Context) error {
  if _, ok := v.touched[key]; !ok {
    return ctx.NewError("Error: this." + key + " not initialized")
  } else {
    return nil
  }
}
