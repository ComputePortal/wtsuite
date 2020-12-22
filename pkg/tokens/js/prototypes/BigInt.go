package prototypes

import (
	"../values"

	"../../context"
)

type BigInt struct {
  BuiltinPrototype
}

func NewBigIntPrototype() values.Prototype {
  return &BigInt{newBuiltinPrototype("BigInt")}
}

func NewBigInt(ctx context.Context) values.Value {
  return values.NewInstance(NewBigIntPrototype(), ctx)
}

func (p *BigInt) GetParent() (values.Prototype, error) {
  return NewIntPrototype(), nil
}

func (p *BigInt) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*BigInt); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *BigInt) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "toLocaleString":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{s, s},
      []values.Value{s, NewLocaleOptions(ctx), s},
    }, ctx), nil
  case "toString":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{i, s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *BigInt) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewClass(
    [][]values.Value{
      []values.Value{NewString(ctx)},
      []values.Value{NewInt(ctx)},
    }, NewBigIntPrototype(), ctx), nil
}
