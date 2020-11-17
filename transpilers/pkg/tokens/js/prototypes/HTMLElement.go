package prototypes

import (
  "../values"

  "../../context"
)

type HTMLElement struct {
  BuiltinPrototype
}

func NewHTMLElementPrototype() values.Prototype {
  return &HTMLElement{newBuiltinPrototype("HTMLElement")}
}

func NewHTMLElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLElementPrototype(), ctx)
}

func (p *HTMLElement) GetParent() (values.Prototype, error) {
  return NewElementPrototype(), nil
}

func (p *HTMLElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)
  elem := NewHTMLElement(ctx)

  switch key {
  case "blur", "click", "focus":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "cellIndex", "rowIndex", "offsetWidth", "offsetHeight":
    return i, nil
  case "children":
    return NewHTMLCollection(ctx), nil
  case "style":
    return NewCSSStyleDeclaration(ctx), nil
  case "parentElement":
    return elem, nil
  case "querySelector":
    return values.NewFunction([]values.Value{s, elem}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *HTMLElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLElementPrototype(), ctx), nil
}
