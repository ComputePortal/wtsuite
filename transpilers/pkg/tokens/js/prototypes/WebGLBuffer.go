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

func (p *WebGLBuffer) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLBufferPrototype(), ctx), nil
}
