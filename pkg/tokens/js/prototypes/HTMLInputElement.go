package prototypes

import (
  "../values"

  "../../context"
)

type HTMLInputElement struct {
  BuiltinPrototype
}

func NewHTMLInputElementPrototype() values.Prototype {
  return &HTMLInputElement{newBuiltinPrototype("HTMLInputElement")}
}

func NewHTMLInputElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLInputElementPrototype(), ctx)
}

func (p *HTMLInputElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLInputElement) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLInputElement); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLInputElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "checked":
    return b, nil
  case "checkValidity":
    return values.NewFunction([]values.Value{b}, ctx), nil
  case "selectionStart", "selectionEnd":
    return i, nil
  case "setCustomValidity":
    return values.NewFunction([]values.Value{s, nil}, ctx), nil
  case "value":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *HTMLInputElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "checked":
    return b.Check(arg, ctx)
  case "selectionStart", "selectionEnd":
    return i.Check(arg, ctx)
  case "value":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLInputElement." + key + " not setable")
  }
}

func (p *HTMLInputElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLInputElementPrototype(), ctx), nil
}
