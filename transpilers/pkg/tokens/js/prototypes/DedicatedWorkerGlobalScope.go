package prototypes

import (
  "../values"

  "../../context"
)

type DedicatedWorkerGlobalScope struct {
  BuiltinPrototype
}

func NewDedicatedWorkerGlobalScopePrototype() values.Prototype {
  return &DedicatedWorkerGlobalScope{newBuiltinPrototype("DedicatedWorkerGlobalScope")}
}

func NewDedicatedWorkerGlobalScope(ctx context.Context) values.Value {
  return values.NewInstance(NewDedicatedWorkerGlobalScopePrototype(), ctx)
}

func NewPostMessageFunction(ctx context.Context) values.Value {
  return values.NewFunction([]values.Value{NewObject(nil, ctx), nil}, ctx)
}

func (p *DedicatedWorkerGlobalScope) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*DedicatedWorkerGlobalScope); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *DedicatedWorkerGlobalScope) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "postMessage":
    return NewPostMessageFunction(ctx), nil
  default:
    return nil, nil
  }
}

func (p *DedicatedWorkerGlobalScope) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewDedicatedWorkerGlobalScopePrototype(), ctx), nil
}
