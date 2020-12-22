package prototypes

import (
  "../values"

  "../../context"
)

type IDBObjectStore struct {
  BuiltinPrototype
}

func NewIDBObjectStorePrototype() values.Prototype {
  return &IDBObjectStore{newBuiltinPrototype("IDBObjectStore")}
}

func NewIDBObjectStore(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBObjectStorePrototype(), ctx)
}

func (p *IDBObjectStore) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBObjectStore); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBObjectStore) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  o := NewObject(nil, ctx)
  s := NewString(ctx)
  blob := NewBlob(ctx)

  switch key {
  case "add":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{o, NewIDBRequest(i, ctx)},
    }, ctx), nil
  case "clear":
    return values.NewFunction([]values.Value{NewEmptyIDBRequest(ctx)}, ctx), nil
  case "count":
    return values.NewFunction([]values.Value{NewIDBRequest(i, ctx)}, ctx), nil
  case "createIndex":
    idx := NewIDBIndex(ctx)
    opt := NewConfigObject(map[string]values.Value{
      "unique": b,
      "multiEntry": b,
      "locale": s,
    }, ctx)

    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{s, s, idx},
      []values.Value{s, s, opt, idx},
    }, ctx), nil
  case "delete":
    req := NewEmptyIDBRequest(ctx)
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, req},
      []values.Value{s, req},
    }, ctx), nil
  case "get":
    req := NewIDBRequest(o, ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, req},
      []values.Value{s, req},
    }, ctx), nil
  case "index":
    idx := NewIDBIndex(ctx)

    return values.NewFunction([]values.Value{s, idx}, ctx), nil
  case "openCursor":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewIDBCursorWithValue(ctx)},
      []values.Value{i, NewIDBCursorWithValue(ctx)},
      []values.Value{NewIDBKeyRange(ctx), NewIDBCursorWithValue(ctx)},
      []values.Value{i, s, NewIDBCursorWithValue(ctx)},
      []values.Value{NewIDBKeyRange(ctx), s, NewIDBCursorWithValue(ctx)},
    }, ctx), nil
  case "put":
    req := NewEmptyIDBRequest(ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{o, req},
      []values.Value{o, i, req},
      []values.Value{o, s, req},
      []values.Value{blob, req},
      []values.Value{blob, i, req},
      []values.Value{blob, s, req},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBObjectStore) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBObjectStorePrototype(), ctx), nil
}
