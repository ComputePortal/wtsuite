package prototypes

import (
  "../values"

  "../../context"
)

type WorkerGlobalScope struct {
  BuiltinPrototype
}

func NewWorkerGlobalScopePrototype() values.Prototype {
  return &WorkerGlobalScope{newBuiltinPrototype("WorkerGlobalScope")}
}

func NewWorkerGlobalScope(ctx context.Context) values.Value {
  return values.NewInstance(NewWorkerGlobalScopePrototype(), ctx)
}

func (p *WorkerGlobalScope) GetParent() (values.Prototype, error) {
  return NewEventTargetPrototype(), nil
}

func (p *WorkerGlobalScope) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWorkerGlobalScopePrototype(), ctx), nil
}
