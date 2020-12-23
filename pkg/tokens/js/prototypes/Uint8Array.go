package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Uint8Array struct {
  AbstractTypedArray
}

func NewUint8ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Uint8Array{newAbstractTypedArrayPrototype("Uint8Array", true, 8, NewInt(ctx))}
}

func NewUint8Array(ctx context.Context) values.Value {
  return values.NewInstance(NewUint8ArrayPrototype(), ctx)
}

func (p *Uint8Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Uint8Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Uint8Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Uint8Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
