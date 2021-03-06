package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type HTMLLinkElement struct {
  BuiltinPrototype
}

func NewHTMLLinkElementPrototype() values.Prototype {
  return &HTMLLinkElement{newBuiltinPrototype("HTMLLinkElement")}
}

func NewHTMLLinkElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLLinkElementPrototype(), ctx)
}

func (p *HTMLLinkElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLLinkElement) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLLinkElement); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLLinkElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "download", "href", "rel":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *HTMLLinkElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  s := NewString(ctx)

  switch key {
  case "download", "href", "rel":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLLinkElement." + key + " not setable")
  }
}

func (p *HTMLLinkElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLLinkElementPrototype(), ctx), nil
}
