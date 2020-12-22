package prototypes

import (
  "../values"

  "../../context"
)

type WebGLShader struct {
  BuiltinPrototype
}

func NewWebGLShaderPrototype() values.Prototype {
  return &WebGLShader{newBuiltinPrototype("WebGLShader")}
}

func NewWebGLShader(ctx context.Context) values.Value {
  return values.NewInstance(NewWebGLShaderPrototype(), ctx)
}

func (p *WebGLShader) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebGLShader); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebGLShader) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLShaderPrototype(), ctx), nil
}
