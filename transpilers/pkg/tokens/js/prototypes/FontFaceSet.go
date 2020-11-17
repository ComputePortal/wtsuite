package prototypes

import (
  "../values"

  "../../context"
)

type FontFaceSet struct {
  BuiltinPrototype
}

func NewFontFaceSetPrototype() values.Prototype {
  return &FontFaceSet{newBuiltinPrototype("FontFaceSet")}
}

func NewFontFaceSet(ctx context.Context) values.Value {
  return values.NewInstance(NewFontFaceSetPrototype(), ctx)
}

func (p *FontFaceSet) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "ready":
    return NewVoidPromise(ctx), nil
  default:
    return nil, nil
  }
}

func (p *FontFaceSet) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewFontFaceSetPrototype(), ctx), nil
}
