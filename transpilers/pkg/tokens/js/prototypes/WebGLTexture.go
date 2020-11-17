package prototypes

import (
  "../values"

  "../../context"
)

type WebGLTexture struct {
  BuiltinPrototype
}

func NewWebGLTexturePrototype() values.Prototype {
  return &WebGLTexture{newBuiltinPrototype("WebGLTexture")}
}

func NewWebGLTexture(ctx context.Context) values.Value {
  return values.NewInstance(NewWebGLTexturePrototype(), ctx)
}

func (p *WebGLTexture) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLTexturePrototype(), ctx), nil
}
