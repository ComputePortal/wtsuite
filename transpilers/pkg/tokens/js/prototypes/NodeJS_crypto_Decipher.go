package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_crypto_Decipher struct {
  BuiltinPrototype
}

func NewNodeJS_crypto_DecipherPrototype() values.Prototype {
  return &NodeJS_crypto_Decipher{newBuiltinPrototype("Decipher")}
}

func NewNodeJS_crypto_Decipher(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_crypto_DecipherPrototype(), ctx)
}

func (p *NodeJS_crypto_Decipher) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_crypto_Decipher) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_crypto_Decipher); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_crypto_Decipher) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "final":
    return values.NewFunction([]values.Value{s, s}, ctx), nil
  case "update":
    return values.NewFunction([]values.Value{s, s, s, s}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_crypto_Decipher) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_crypto_DecipherPrototype(), ctx), nil
}
