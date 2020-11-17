package prototypes

import (
  "../values"

  "../../context"
)

type Uint8ClampedArray struct {
  AbstractTypedArray
}

func NewUint8ClampedArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Uint8ClampedArray{newAbstractTypedArrayPrototype("Uint8ClampedArray", true, 8, NewInt(ctx))}
}

func NewUint8ClampedArray(ctx context.Context) values.Value {
  return values.NewInstance(NewUint8ClampedArrayPrototype(), ctx)
}

func (p *Uint8ClampedArray) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Uint8ClampedArray) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Uint8ClampedArray) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Uint8ClampedArray) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
