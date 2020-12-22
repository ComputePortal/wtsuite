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

func (p *WebGLTexture) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebGLTexture); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebGLTexture) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLTexturePrototype(), ctx), nil
}
