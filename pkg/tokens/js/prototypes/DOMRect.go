package prototypes

import (
  "../values"

  "../../context"
)

type DOMRect struct {
  BuiltinPrototype
}

func NewDOMRectPrototype() values.Prototype {
  return &DOMRect{newBuiltinPrototype("DOMRect")}
}

func NewDOMRect(ctx context.Context) values.Value {
  return values.NewInstance(NewDOMRectPrototype(), ctx)
}

func (p *DOMRect) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*DOMRect); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *DOMRect) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  f := NewNumber(ctx)

  switch key {
  case "bottom", "height", "left", "right", "top", "width", "x", "y":
    return f, nil
  default:
    return nil, nil
  }
}

func (p *DOMRect) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewDOMRectPrototype(), ctx), nil
}
