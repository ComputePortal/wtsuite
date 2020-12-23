package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Int8Array struct {
  AbstractTypedArray
}

func NewInt8ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Int8Array{newAbstractTypedArrayPrototype("Int8Array", false, 8, NewInt(ctx))}
}

func NewInt8Array(ctx context.Context) values.Value {
  return values.NewInstance(NewInt8ArrayPrototype(), ctx)
}

func (p *Int8Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Int8Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Int8Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Int8Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
