package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Uint16Array struct {
  AbstractTypedArray
}

func NewUint16ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Uint16Array{newAbstractTypedArrayPrototype("Uint16Array", true, 16, NewInt(ctx))}
}

func NewUint16Array(ctx context.Context) values.Value {
  return values.NewInstance(NewUint16ArrayPrototype(), ctx)
}

func (p *Uint16Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Uint16Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Uint16Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Uint16Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
