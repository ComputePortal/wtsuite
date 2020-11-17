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

func (p *DOMRect) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  s := NewString(ctx)

  switch key {
  case "bottom", "height", "left", "right", "top", "width", "x", "y":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: DOMRect." + key + " not setable")
  }
}

func (p *DOMRect) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewDOMRectPrototype(), ctx), nil
}
