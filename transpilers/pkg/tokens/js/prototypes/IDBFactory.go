package prototypes

import (
  "../values"

  "../../context"
)

type IDBFactory struct {
  BuiltinPrototype
}

func NewIDBFactoryPrototype() values.Prototype {
  return &IDBFactory{newBuiltinPrototype("IDBFactory")}
}

func NewIDBFactory(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBFactoryPrototype(), ctx)
}

func (p *IDBFactory) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBFactory); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBFactory) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)
  req := NewIDBOpenDBRequest(ctx)

  switch key {
  case "databases":
    return values.NewFunction([]values.Value{
      NewPromise(NewArray(NewObject(map[string]values.Value{
        "name": s,
        "version": i,
      }, ctx), ctx), ctx),
    }, ctx), nil
  case "deleteDatabase":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, req},
      // []values.Value{s, o, req}, // not yet standardized
    }, ctx), nil
  case "open":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, req},
      []values.Value{s, i, req},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBFactory) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBFactoryPrototype(), ctx), nil
}
