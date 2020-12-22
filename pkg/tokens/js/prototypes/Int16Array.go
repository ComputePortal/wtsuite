package prototypes

import (
  "../values"

  "../../context"
)

type Int16Array struct {
  AbstractTypedArray
}

func NewInt16ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Int16Array{newAbstractTypedArrayPrototype("Int16Array", false, 16, NewInt(ctx))}
}

func NewInt16Array(ctx context.Context) values.Value {
  return values.NewInstance(NewInt16ArrayPrototype(), ctx)
}

func (p *Int16Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Int16Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Int16Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Int16Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
