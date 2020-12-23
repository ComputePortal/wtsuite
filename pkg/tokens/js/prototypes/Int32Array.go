package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Int32Array struct {
  AbstractTypedArray
}

func NewInt32ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Int32Array{newAbstractTypedArrayPrototype("Int32Array", false, 32, NewInt(ctx))}
}

func NewInt32Array(ctx context.Context) values.Value {
  return values.NewInstance(NewInt32ArrayPrototype(), ctx)
}

func (p *Int32Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Int32Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Int32Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Int32Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
