package prototypes

import (
  "../values"

  "../../context"
)

type Response struct {
  BuiltinPrototype
}

func NewResponsePrototype() values.Prototype {
  return &Response{newBuiltinPrototype("Response")}
}

func NewResponse(ctx context.Context) values.Value {
  return values.NewInstance(NewResponsePrototype(), ctx)
}

func (p *Response) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "blob":
    return values.NewFunction([]values.Value{NewPromise(NewBlob(ctx), ctx)}, ctx), nil
  case "json":
    return values.NewFunction([]values.Value{NewPromise(NewObject(nil, ctx), ctx)}, ctx), nil
  case "ok":
    return b, nil
  case "status":
    return i, nil
  case "statusText":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *Response) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewResponsePrototype(), ctx), nil
}
