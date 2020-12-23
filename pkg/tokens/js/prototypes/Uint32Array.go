package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Uint32Array struct {
  AbstractTypedArray
}

func NewUint32ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Uint32Array{newAbstractTypedArrayPrototype("Uint32Array", true, 32, NewInt(ctx))}
}

func NewUint32Array(ctx context.Context) values.Value {
  return values.NewInstance(NewUint32ArrayPrototype(), ctx)
}

func (p *Uint32Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Uint32Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Uint32Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Uint32Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
