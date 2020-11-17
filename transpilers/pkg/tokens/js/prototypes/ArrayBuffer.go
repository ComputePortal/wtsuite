package prototypes

import (
  "../values"

  "../../context"
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

func (p *ArrayBuffer) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewArrayBufferPrototype(), ctx), nil
}
