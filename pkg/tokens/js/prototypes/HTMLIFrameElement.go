package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type HTMLIFrameElement struct {
  BuiltinPrototype
}

func NewHTMLIFrameElementPrototype() values.Prototype {
  return &HTMLIFrameElement{newBuiltinPrototype("HTMLIFrameElement")}
}

func NewHTMLIFrameElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLIFrameElementPrototype(), ctx)
}

func (p *HTMLIFrameElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLIFrameElement) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLIFrameElement); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLIFrameElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "contentDocument":
    return NewDocument(ctx), nil
  case "contentWindow":
    return NewWindow(ctx), nil
  case "src":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *HTMLIFrameElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  s := NewString(ctx)

  switch key {
  case "src":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLIFrameElement." + key + " not setable")
  }
}

func (p *HTMLIFrameElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLIFrameElementPrototype(), ctx), nil
}
