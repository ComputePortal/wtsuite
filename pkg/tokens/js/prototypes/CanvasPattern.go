package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
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

func (p *CanvasPattern) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*CanvasPattern); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *CanvasPattern) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCanvasPatternPrototype(), ctx), nil
}
