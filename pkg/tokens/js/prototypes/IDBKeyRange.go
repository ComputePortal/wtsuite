package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type IDBKeyRange struct {
  BuiltinPrototype
}

func NewIDBKeyRangePrototype() values.Prototype {
  return &IDBKeyRange{newBuiltinPrototype("IDBKeyRange")}
}

func NewIDBKeyRange(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBKeyRangePrototype(), ctx)
}

func (p *IDBKeyRange) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBKeyRange); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBKeyRange) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)

  switch key {
  case "includes":
    return values.NewFunction([]values.Value{i, b}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBKeyRange) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)
  self := NewIDBKeyRange(ctx)

  switch key {
  case "lowerBound", "upperBound":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, self},
      []values.Value{i, b, self},
      []values.Value{s, self},
      []values.Value{s, b, self},
    }, ctx), nil
  case "bound":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, i, self},
      []values.Value{i, i, b, self},
      []values.Value{i, i, b, b, self},
      []values.Value{s, s, self},
      []values.Value{s, s, b, self},
      []values.Value{s, s, b, b, self},
    }, ctx), nil
  case "only":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, self},
      []values.Value{s, self},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *IDBKeyRange) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBKeyRangePrototype(), ctx), nil
}
