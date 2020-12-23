package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Console struct {
  BuiltinPrototype
}

func NewConsolePrototype() values.Prototype {
  return &Console{newBuiltinPrototype("Console")}
}

func NewConsole(ctx context.Context) values.Value {
  return values.NewInstance(NewConsolePrototype(), ctx)
}

func (p *Console) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Console); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Console) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)

  switch key {
  case "log":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{nil},
      []values.Value{a, nil},
      []values.Value{a, a, nil},
      []values.Value{a, a, a, nil},
      []values.Value{a, a, a, a, nil},
      []values.Value{a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, nil}, // should be enough
    }, ctx), nil
  case "error":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{nil},
      []values.Value{a, nil},
      []values.Value{a, a, nil},
      []values.Value{a, a, a, nil},
      []values.Value{a, a, a, a, nil},
      []values.Value{a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, a, a, nil},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, nil}, // should be enough
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Console) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewConsolePrototype(), ctx), nil
}
