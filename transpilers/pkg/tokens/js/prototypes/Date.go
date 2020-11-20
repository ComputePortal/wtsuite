package prototypes

import (
  "../values"

  "../../context"
)

type Date struct {
  BuiltinPrototype
}

func NewDatePrototype() values.Prototype {
  return &Date{newBuiltinPrototype("Date")}
}

func NewDate(ctx context.Context) values.Value {
  return values.NewInstance(NewDatePrototype(), ctx)
}

func (p *Date) IsUniversal() bool {
  return true
}

func (p *Date) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Date); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Date) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "getTime":
    return values.NewFunction([]values.Value{i}, ctx), nil
  case "setTime":
    return values.NewMethodLikeFunction([]values.Value{i, i}, ctx), nil
  case "toGMTString":
    return values.NewFunction([]values.Value{s}, ctx), nil
  case "toLocaleString":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{s, s},
      []values.Value{s, NewLocaleOptions(ctx), s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Date) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  self := NewDate(ctx)

  switch key {
  case "now":
    return values.NewFunction([]values.Value{self}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Date) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewClass([][]values.Value{
    []values.Value{},
    []values.Value{NewNumber(ctx)},
    []values.Value{NewString(ctx)},
  }, NewDatePrototype(), ctx), nil
}
