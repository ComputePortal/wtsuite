package prototypes

import (
  "../values"

  "../../context"
)

type CSSImageValue struct {
  BuiltinPrototype
}

func NewCSSImageValuePrototype() values.Prototype {
  return &CSSImageValue{newBuiltinPrototype("CSSImageValue")}
}

func NewCSSImageValue(ctx context.Context) values.Value {
  return values.NewInstance(NewCSSImageValuePrototype(), ctx)
}

func (p *CSSImageValue) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*CSSImageValue); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *CSSImageValue) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCSSImageValuePrototype(), ctx), nil
}
