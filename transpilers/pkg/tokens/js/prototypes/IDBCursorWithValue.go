package prototypes

import (
  "../values"

  "../../context"
)

type IDBCursorWithValue struct {
  BuiltinPrototype
}

func NewIDBCursorWithValuePrototype() values.Prototype {
  return &IDBCursorWithValue{newBuiltinPrototype("IDBCursorWithValue")}
}

func NewIDBCursorWithValue(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBCursorWithValuePrototype(), ctx)
}

func (p *IDBCursorWithValue) GetParent() (values.Prototype, error) {
  return NewIDBCursorPrototype(), nil
}

func (p *IDBCursorWithValue) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "value":
    return NewObject(nil, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBCursorWithValue) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBCursorWithValuePrototype(), ctx), nil
}
