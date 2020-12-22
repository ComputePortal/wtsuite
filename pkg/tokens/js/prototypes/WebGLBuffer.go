package prototypes

import (
  "../values"

  "../../context"
)

type WebGLBuffer struct {
  BuiltinPrototype
}

func NewWebGLBufferPrototype() values.Prototype {
  return &WebGLBuffer{newBuiltinPrototype("WebGLBuffer")}
}

func NewWebGLBuffer(ctx context.Context) values.Value {
  return values.NewInstance(NewWebGLBufferPrototype(), ctx)
}

func (p *WebGLBuffer) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebGLBuffer); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebGLBuffer) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLBufferPrototype(), ctx), nil
}
