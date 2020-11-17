package prototypes

import (
  "../values"

  "../../context"
)

type JSON struct {
  BuiltinPrototype
}

func NewJSONPrototype() values.Prototype {
  return &JSON{newBuiltinPrototype("JSON")}
}

func NewJSON(ctx context.Context) values.Value {
  return values.NewInstance(NewJSONPrototype(), ctx)
}

func (p *JSON) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)
  o := NewObject(nil, ctx)

  switch key {
  case "stringify":
    return values.NewFunction([]values.Value{o, s}, ctx), nil
  case "parse":
    return values.NewFunction([]values.Value{s, o}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *JSON) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewJSONPrototype(), ctx), nil
}
