package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ArrayBuffer struct {
  BuiltinPrototype
}

func NewArrayBufferPrototype() values.Prototype {
  return &ArrayBuffer{newBuiltinPrototype("ArrayBuffer")}
}

func NewArrayBuffer(ctx context.Context) values.Value {
  return values.NewInstance(NewArrayBufferPrototype(), ctx)
}

func (p *ArrayBuffer) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*ArrayBuffer); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *ArrayBuffer) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewArrayBufferPrototype(), ctx), nil
}
