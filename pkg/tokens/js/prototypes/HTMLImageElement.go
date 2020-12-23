package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type HTMLImageElement struct {
  BuiltinPrototype
}

func NewHTMLImageElementPrototype() values.Prototype {
  return &HTMLImageElement{newBuiltinPrototype("HTMLImageElement")}
}

func NewHTMLImageElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLImageElementPrototype(), ctx)
}

func (p *HTMLImageElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLImageElement) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLImageElement); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLImageElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)

  switch key {
  case "onload", "src":
    return nil, ctx.NewError("Error: is only a setter")
  case "height", "width":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *HTMLImageElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  i := NewInt(ctx)
  s := NewString(ctx)
  self := values.NewInstance(p, ctx)

  switch key {
  case "height", "width":
    return i.Check(arg, ctx)
  case "onload":
    callback := values.NewFunction([]values.Value{NewEvent(self, ctx), nil}, ctx)

    return callback.Check(arg, ctx)
  case "src":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLImageElement." + key + " not setable")
  }
}

func (p *HTMLImageElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLImageElementPrototype(), ctx), nil
}
