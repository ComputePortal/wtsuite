package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Float32Array struct {
  AbstractTypedArray
}

func NewFloat32ArrayPrototype() values.Prototype {
  ctx := context.NewDummyContext()

  return &Float32Array{newAbstractTypedArrayPrototype("Float32Array", false, 32, NewNumber(ctx))}
}

func NewFloat32Array(ctx context.Context) values.Value {
  return values.NewInstance(NewFloat32ArrayPrototype(), ctx)
}

func (p *Float32Array) Check(other values.Interface, ctx context.Context) error {
  return CheckTypedArray(p, other, ctx)
}

func (p *Float32Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayInstanceMember(p, key, includePrivate, ctx)
}

func (p *Float32Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return GetTypedArrayClassMember(p, key, includePrivate, ctx)
}

func (p *Float32Array) GetClassValue() (*values.Class, error) {
  return GetTypedArrayClassValue(p)
}
