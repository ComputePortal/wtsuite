package prototypes

import (
  "../values"

  "../../context"
)

type HTMLTextAreaElement struct {
  BuiltinPrototype
}

func NewHTMLTextAreaElementPrototype() values.Prototype {
  return &HTMLTextAreaElement{newBuiltinPrototype("HTMLTextAreaElement")}
}

func NewHTMLTextAreaElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLTextAreaElementPrototype(), ctx)
}

func (p *HTMLTextAreaElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLTextAreaElement) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLTextAreaElement); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLTextAreaElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
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

func (p *HTMLTextAreaElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "selectionStart", "selectionEnd":
    return i.Check(arg, ctx)
  case "value":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLTextAreaElement." + key + " not setable")
  }
}

func (p *HTMLTextAreaElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLTextAreaElementPrototype(), ctx), nil
}
