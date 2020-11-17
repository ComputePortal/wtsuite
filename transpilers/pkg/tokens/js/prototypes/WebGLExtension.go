package prototypes

import (
  "../values"

  "../../context"
)

type WebGLExtension struct {
  BuiltinPrototype
}

func NewWebGLExtensionPrototype() values.Prototype {
  return &WebGLExtension{newBuiltinPrototype("WebGLExtension")}
}

func NewWebGLExtension(ctx context.Context) values.Value {
  return values.NewInstance(NewWebGLExtensionPrototype(), ctx)
}

func (p *WebGLExtension) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLExtensionPrototype(), ctx), nil
}
