package prototypes

import (
  "../values"

  "../../context"
)

type RegExpArray struct {
  BuiltinPrototype
}

func NewRegExpArrayPrototype() values.Prototype {
  return &RegExpArray{newBuiltinPrototype("RegExpArray")}
}

func NewRegExpArray(ctx context.Context) values.Value {
  return values.NewInstance(NewRegExpArrayPrototype(), ctx)
}

func (p *RegExpArray) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*RegExpArray); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *RegExpArray) GetParent() (values.Prototype, error) {
  return NewArrayPrototype(NewString(p.Context())), nil
}

func (p *RegExpArray) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "index":
    return i, nil
  case "input":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *RegExpArray) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewUnconstructableClass(NewRegExpArrayPrototype(), ctx), nil
}
