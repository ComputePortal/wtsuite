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

func (p *CSSImageValue) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCSSImageValuePrototype(), ctx), nil
}
