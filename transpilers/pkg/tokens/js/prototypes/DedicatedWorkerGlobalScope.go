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
  a := values.NewAny(ctx)

  return values.NewFunction([]values.Value{a, nil}, ctx)
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
