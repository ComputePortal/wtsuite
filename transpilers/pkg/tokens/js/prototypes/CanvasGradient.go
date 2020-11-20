package prototypes

import (
  "../values"

  "../../context"
)

type CanvasGradient struct {
  BuiltinPrototype
}

func NewCanvasGradientPrototype() values.Prototype {
  return &CanvasGradient{newBuiltinPrototype("CanvasGradient")}
}

func NewCanvasGradient(ctx context.Context) values.Value {
  return values.NewInstance(NewCanvasGradientPrototype(), ctx)
}

func IsCanvasGradient(v values.Value) bool {
  ctx := v.Context()

  checkVal := NewCanvasGradient(ctx)

  return checkVal.Check(v, ctx) == nil
}

func (p *CanvasGradient) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*CanvasGradient); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *CanvasGradient) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "addColorStop":
    return values.NewFunction([]values.Value{NewNumber(ctx), NewString(ctx), nil}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *CanvasGradient) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCanvasGradientPrototype(), ctx), nil
}
