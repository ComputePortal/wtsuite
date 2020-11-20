package prototypes

import (
  "../values"

  "../../context"
)

type IDBDatabase struct {
  BuiltinPrototype
}

func NewIDBDatabasePrototype() values.Prototype {
  return &IDBDatabase{newBuiltinPrototype("IDBDatabase")}
}

func NewIDBDatabase(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBDatabasePrototype(), ctx)
}

func (p *IDBDatabase) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBDatabase); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBDatabase) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  s := NewString(ctx)
  ss := NewArray(s, ctx)
  store := NewIDBObjectStore(ctx)
  trans := NewIDBTransaction(ctx)

  switch key {
  case "close":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "createObjectStore":
    o := NewConfigObject(map[string]values.Value{
      "keyPath": s,
      "autoIncrement": b,
    }, ctx)

    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{s, store},
      []values.Value{s, o, store},
    }, ctx), nil
  case "name":
    return s, nil
  case "transaction":
    opt := NewConfigObject(map[string]values.Value{
      "durability": s,
    }, ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, trans},
      []values.Value{ss, trans},
      []values.Value{s, s, trans},
      []values.Value{ss, s, trans},
      []values.Value{s, s, opt, trans},
      []values.Value{ss, s, opt, trans},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBDatabase) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBDatabasePrototype(), ctx), nil
}
