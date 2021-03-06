package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type IDBTransaction struct {
  BuiltinPrototype
}

func NewIDBTransactionPrototype() values.Prototype {
  return &IDBTransaction{newBuiltinPrototype("IDBTransaction")}
}

func NewIDBTransaction(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBTransactionPrototype(), ctx)
}

func (p *IDBTransaction) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBTransaction); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBTransaction) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)
  store := NewIDBObjectStore(ctx)

  switch key {
  case "oncomplete", "onerror":
    return nil, ctx.NewError("Error: is only a setter")
  case "commit":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "objectStore":
    return values.NewFunction([]values.Value{s, store}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBTransaction) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  callback := values.NewFunction([]values.Value{NewEvent(NewIDBTransaction(ctx), ctx), nil}, ctx)

  switch key {
  case "oncomplete", "onerror":
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: IDBTransaction." + key + " not setable")
  }
}

func (p *IDBTransaction) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBTransactionPrototype(), ctx), nil
}
