package prototypes

import (
  "../values"

  "../../context"
)

type Image struct {
  BuiltinPrototype
}

func NewImagePrototype() values.Prototype {
  return &Image{newBuiltinPrototype("Image")}
}

func NewImage(ctx context.Context) values.Value {
  return values.NewInstance(NewImagePrototype(), ctx)
}

func (p *Image) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Image); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Image) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  f := NewNumber(ctx)

  // TODO: special class values type that display Image as name, but result of constructor is HTMLImageElement
  return values.NewClass([][]values.Value{
    []values.Value{f, f},
  }, NewHTMLImageElementPrototype(), ctx), nil
}
