package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
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

func (p *SharedWorkerGlobalScope) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*SharedWorkerGlobalScope); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *SharedWorkerGlobalScope) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewSharedWorkerGlobalScopePrototype(), ctx), nil
}
