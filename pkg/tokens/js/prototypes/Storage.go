package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Storage struct {
  BuiltinPrototype
}

func NewStoragePrototype() values.Prototype {
  return &Storage{newBuiltinPrototype("Storage")}
}

func NewStorage(ctx context.Context) values.Value {
  return values.NewInstance(NewStoragePrototype(), ctx)
}

func (p *Storage) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Storage); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Storage) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "clear":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "getItem":
    return values.NewFunction([]values.Value{s, s}, ctx), nil
  case "key":
    return values.NewFunction([]values.Value{i, s}, ctx), nil
  case "length":
    return i, nil
  case "removeItem":
    return values.NewFunction([]values.Value{s, nil}, ctx), nil
  case "setItem":
    return values.NewFunction([]values.Value{s, s, nil}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Storage) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewStoragePrototype(), ctx), nil
}
