package prototypes

import (
  "../values"

  "../../context"
)

type IDBIndex struct {
  BuiltinPrototype
}

func NewIDBIndexPrototype() values.Prototype {
  return &IDBIndex{newBuiltinPrototype("IDBIndex")}
}

func NewIDBIndex(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBIndexPrototype(), ctx)
}

func (p *IDBIndex) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBIndex); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBIndex) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  kr := NewIDBKeyRange(ctx)
  s := NewString(ctx)

  switch key {
  case "getAll":
    req := NewIDBRequest(NewArray(NewObject(nil, ctx), ctx), ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{req},
      []values.Value{kr, req},
      []values.Value{i, req},
    }, ctx), nil
  case "getAllKeys":
    req := NewIDBRequest(NewArray(i, ctx), ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{req},
      []values.Value{kr, req},
      []values.Value{i, req},
    }, ctx), nil
  case "openCursor":
    req := NewIDBRequest(NewIDBCursorWithValue(ctx), ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{req},
      []values.Value{i, req},
      []values.Value{kr, req},
      []values.Value{i, s, req},
      []values.Value{kr, s, req},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBIndex) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBIndexPrototype(), ctx), nil
}
