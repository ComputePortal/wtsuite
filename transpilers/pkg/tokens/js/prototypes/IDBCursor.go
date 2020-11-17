package prototypes

import (
  "../values"

  "../../context"
)

type IDBCursor struct {
  BuiltinPrototype
}

func NewIDBCursorPrototype() values.Prototype {
  return &IDBCursor{newBuiltinPrototype("IDBCursor")}
}

func NewIDBCursor(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBCursorPrototype(), ctx)
}

func (p *IDBCursor) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)

  switch key {
  case "advance":
    return values.NewFunction([]values.Value{i, nil}, ctx), nil
  case "continue":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{nil}, 
      []values.Value{i, nil},
    }, ctx), nil
  case "continuePrimaryKey":
    return values.NewFunction([]values.Value{i, i, nil}, ctx), nil
  case "delete":
    return values.NewFunction([]values.Value{NewEmptyIDBRequest(ctx)}, ctx), nil
  case "key":
    return i, nil
  case "update":
    return values.NewFunction([]values.Value{NewObject(nil, ctx), NewEmptyIDBRequest(ctx)}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBCursor) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBCursorPrototype(), ctx), nil
}
