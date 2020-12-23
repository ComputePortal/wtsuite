package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type WebAssemblyEnv struct {
  BuiltinPrototype
}

func NewWebAssemblyEnvPrototype() values.Prototype {
  return &WebAssemblyEnv{newBuiltinPrototype("WebAssemblyEnv")}
}

func NewWebAssemblyEnv(ctx context.Context) values.Value {
  return values.NewInstance(NewWebAssemblyEnvPrototype(), ctx)
}

func IsWebAssemblyEnv(v values.Value) bool {
  ctx := context.NewDummyContext()

  checkVal := NewWebAssemblyEnv(ctx)

  return checkVal.Check(v, ctx) == nil
}

func (p *WebAssemblyEnv) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebAssemblyEnv); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebAssemblyEnv) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewClass([][]values.Value{
    []values.Value{NewWebAssemblyFS(ctx)},
  }, NewWebAssemblyEnvPrototype(), ctx), nil
}
