package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
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

func (p *WebGLExtension) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebGLExtension); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebGLExtension) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLExtensionPrototype(), ctx), nil
}
