package prototypes

import (
  "../values"

  "../../context"
)

type Path2D struct {
  BuiltinPrototype
}

func NewPath2DPrototype() values.Prototype {
  return &Path2D{newBuiltinPrototype("Path2D")}
}

func NewPath2D(ctx context.Context) values.Value {
  return values.NewInstance(NewPath2DPrototype(), ctx)
}

func (p *Path2D) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Path2D); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Path2D) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewPath2DPrototype(), ctx), nil
}
