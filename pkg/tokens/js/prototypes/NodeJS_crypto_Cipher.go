package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type NodeJS_crypto_Cipher struct {
  BuiltinPrototype
}

func NewNodeJS_crypto_CipherPrototype() values.Prototype {
  return &NodeJS_crypto_Cipher{newBuiltinPrototype("Cipher")}
}

func NewNodeJS_crypto_Cipher(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_crypto_CipherPrototype(), ctx)
}

func (p *NodeJS_crypto_Cipher) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_crypto_Cipher) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_crypto_Cipher); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_crypto_Cipher) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
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

func (p *NodeJS_crypto_Cipher) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()
  return values.NewUnconstructableClass(NewNodeJS_crypto_CipherPrototype(), ctx), nil
}
