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

func (p *WorkerGlobalScope) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Element); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WorkerGlobalScope) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWorkerGlobalScopePrototype(), ctx), nil
}
