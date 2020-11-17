package prototypes

import (
  "../values"

  "../../context"
)

type SharedWorkerGlobalScope struct {
  BuiltinPrototype
}

func NewSharedWorkerGlobalScopePrototype() values.Prototype {
  return &SharedWorkerGlobalScope{newBuiltinPrototype("SharedWorkerGlobalScope")}
}

func NewSharedWorkerGlobalScope(ctx context.Context) values.Value {
  return values.NewInstance(NewSharedWorkerGlobalScopePrototype(), ctx)
}

func (p *SharedWorkerGlobalScope) GetParent() (values.Prototype, error) {
  return NewWorkerGlobalScopePrototype(), nil
}

func (p *SharedWorkerGlobalScope) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewSharedWorkerGlobalScopePrototype(), ctx), nil
}
