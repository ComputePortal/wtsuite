package prototypes

import (
  "../values"

  "../../context"
)

type HTMLSelectElement struct {
  BuiltinPrototype
}

func NewHTMLSelectElementPrototype() values.Prototype {
  return &HTMLSelectElement{newBuiltinPrototype("HTMLSelectElement")}
}

func NewHTMLSelectElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLSelectElementPrototype(), ctx)
}

func (p *HTMLSelectElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLSelectElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "checkValidity":
    return values.NewFunction([]values.Value{b}, ctx), nil
  case "selectedIndex":
    return i, nil
  case "setCustomValidity":
    return values.NewFunction([]values.Value{s, nil}, ctx), nil
  case "value":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *HTMLSelectElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "selectedIndex":
    return i.Check(arg, ctx)
  case "value":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLSelectElement." + key + " not setable")
  }
}

func (p *HTMLSelectElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLSelectElementPrototype(), ctx), nil
}
