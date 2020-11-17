package prototypes

import (
  "../values"

  "../../context"
)

type CanvasPattern struct {
  BuiltinPrototype
}

func NewCanvasPatternPrototype() values.Prototype {
  return &CanvasPattern{newBuiltinPrototype("CanvasPattern")}
}

func NewCanvasPattern(ctx context.Context) values.Value {
  return values.NewInstance(NewCanvasPatternPrototype(), ctx)
}

func IsCanvasPattern(v values.Value) bool {
  ctx := v.Context()

  checkVal := NewCanvasPattern(ctx)

  return checkVal.Check(v, ctx) == nil
}

func (p *CanvasPattern) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCanvasPatternPrototype(), ctx), nil
}
