package prototypes

import (
  "../values"

  "../../context"
)

type TextEncoder struct {
  BuiltinPrototype
}

func NewTextEncoderPrototype() values.Prototype {
  return &TextEncoder{newBuiltinPrototype("TextEncoder")}
}

func NewTextEncoder(ctx context.Context) values.Value {
  return values.NewInstance(NewTextEncoderPrototype(), ctx)
}

func (p *TextEncoder) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "encode": 
    return values.NewFunction([]values.Value{s, NewUint8Array(ctx)}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *TextEncoder) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewClass([][]values.Value{
    []values.Value{},
  }, NewTextEncoderPrototype(), ctx), nil
}
