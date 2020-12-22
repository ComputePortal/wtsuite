package prototypes

import (
  "../values"

  "../../context"
)

type WebAssembly struct {
  BuiltinPrototype
}

func NewWebAssemblyPrototype() values.Prototype {
  return &WebAssembly{newBuiltinPrototype("WebAssembly")}
}

func NewWebAssembly(ctx context.Context) values.Value {
  return values.NewInstance(NewWebAssemblyPrototype(), ctx)
}

func (p *WebAssembly) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebAssembly); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebAssembly) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebAssemblyPrototype(), ctx), nil
}
