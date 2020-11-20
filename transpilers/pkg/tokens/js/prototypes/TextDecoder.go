package prototypes

import (
  "../values"

  "../../context"
)

type TextDecoder struct {
  BuiltinPrototype
}

func NewTextDecoderPrototype() values.Prototype {
  return &TextDecoder{newBuiltinPrototype("TextDecoder")}
}

func NewTextDecoder(ctx context.Context) values.Value {
  return values.NewInstance(NewTextDecoderPrototype(), ctx)
}

func (p *TextDecoder) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*TextDecoder); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *TextDecoder) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "decode": 
    return values.NewFunction([]values.Value{NewUint8Array(ctx), s}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *TextDecoder) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  b := NewBoolean(ctx)
  s := NewString(ctx)

  return values.NewClass([][]values.Value{
    []values.Value{},
    []values.Value{s},
    []values.Value{s, b},
  }, NewTextDecoderPrototype(), ctx), nil
}
