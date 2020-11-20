package prototypes

import (
  "../values"

  "../../context"
)

type ImageData struct {
  BuiltinPrototype
}

func NewImageDataPrototype() values.Prototype {
  return &ImageData{newBuiltinPrototype("ImageData")}
}

func NewImageData(ctx context.Context) values.Value {
  return values.NewInstance(NewImageDataPrototype(), ctx)
}

func (p *ImageData) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*ImageData); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *ImageData) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)

  switch key {
  case "data":
    return NewUint8ClampedArray(ctx), nil
  case "height", "width":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *ImageData) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewImageDataPrototype(), ctx), nil
}
