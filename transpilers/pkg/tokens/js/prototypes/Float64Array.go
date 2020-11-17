package prototypes

import (
  "../values"

  "../../context"
)

type Float64Array struct {
  AbstractTypedArray
}

func NewFloat64ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Float64Array{newAbstractTypedArrayPrototype("Float64Array", false, 64, NewNumber(ctx))}
}

func NewFloat64Array(ctx context.Context) values.Value {
  return values.NewInstance(NewFloat64ArrayPrototype(), ctx)
}

func (p *Float64Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Float64Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Float64Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Float64Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
