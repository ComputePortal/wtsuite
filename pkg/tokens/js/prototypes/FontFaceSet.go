package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
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

func (p *FontFaceSet) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*FontFaceSet); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
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
