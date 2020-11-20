package values

import (
  "../../context"
)

// used in constructor, to allow calling
type Super struct {
  cl    *Class
  super Value
  ValueData
}

func NewSuper(cl *Class, ctx context.Context) (Value, error) {
  if cl == nil {
    return nil, ctx.NewError("Error: super not found")
  }

  val, err := cl.EvalConstructor(nil, ctx)
  if err != nil {
    return nil, err
  }

  return &Super{
    cl,
    val,
    ValueData{ctx},
  }, nil
}

func (v *Super) TypeName() string {
  return v.super.TypeName()
}

func (v *Super) Check(other_ Value, ctx context.Context) error {
  return v.super.Check(other_, ctx)
}

func (v *Super) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't construct super")
}

// always acts as method
func (v *Super) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  if  _, err := v.cl.EvalConstructor(args, ctx); err != nil {
    return nil, err
  }

  return nil, nil
}

func (v *Super) GetMember(key string, includePrivate bool, ctx context.Context) (Value, error) {
  return v.super.GetMember(key, includePrivate, ctx)
}

func (v *Super) SetMember(key string, includePrivate bool, arg Value, ctx context.Context) error {
  return v.super.SetMember(key, includePrivate, arg, ctx)
}
