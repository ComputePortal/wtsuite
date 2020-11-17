package prototypes

import (
  "../values"

  "../../context"
)

type URLSearchParams struct {
  BuiltinPrototype
}

func NewURLSearchParamsPrototype() values.Prototype {
  return &URLSearchParams{newBuiltinPrototype("URLSearchParams")}
}

func NewURLSearchParams(ctx context.Context) values.Value {
  return values.NewInstance(NewURLSearchParamsPrototype(), ctx)
}

func (p *URLSearchParams) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  s := NewString(ctx)

  switch key {
  case "get":
    return values.NewFunction([]values.Value{s, s}, ctx), nil
  case "has":
    return values.NewFunction([]values.Value{s, b}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *URLSearchParams) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewURLSearchParamsPrototype(), ctx), nil
}
