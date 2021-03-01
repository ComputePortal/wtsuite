package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Crypto struct {
  BuiltinPrototype
}

func NewCryptoPrototype() values.Prototype {
  return &Crypto{newBuiltinPrototype("Crypto")}
}

func NewCrypto(ctx context.Context) values.Value {
  return values.NewInstance(NewCryptoPrototype(), ctx)
}

func (p *Crypto) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Crypto); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Crypto) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := NewTypedArray(ctx)

  switch key {
  case "getRandomValues":
    return values.NewMethodLikeFunction([]values.Value{a, a}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Crypto) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCryptoPrototype(), ctx), nil
}
