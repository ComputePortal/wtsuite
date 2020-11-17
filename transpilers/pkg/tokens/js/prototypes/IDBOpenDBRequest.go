package prototypes

import (
  "../values"

  "../../context"
)

type IDBOpenDBRequest struct {
  BuiltinPrototype
}

func NewIDBOpenDBRequestPrototype() values.Prototype {
  return &IDBOpenDBRequest{newBuiltinPrototype("IDBOpenDBRequest")}
}

func NewIDBOpenDBRequest(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBOpenDBRequestPrototype(), ctx)
}

func (p *IDBOpenDBRequest) GetParent() (values.Prototype, error) {
  ctx := p.Context()
  return NewIDBRequestPrototype(NewIDBDatabase(ctx)), nil
}

func (p *IDBOpenDBRequest) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  callback := values.NewFunction([]values.Value{NewIDBVersionChangeEvent(ctx)}, ctx)

  switch key {
  case "onupgradeneeded":
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: IDBOpenDBRequest." + key + " not setable")
  }
}

func (p *IDBOpenDBRequest) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBOpenDBRequestPrototype(), ctx), nil
}
