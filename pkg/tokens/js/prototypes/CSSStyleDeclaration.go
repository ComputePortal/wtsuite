package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type CSSStyleDeclaration struct {
  BuiltinPrototype
}

func NewCSSStyleDeclarationPrototype() values.Prototype {
  return &CSSStyleDeclaration{newBuiltinPrototype("CSSStyleDeclaration")}
}

func NewCSSStyleDeclaration(ctx context.Context) values.Value {
  return values.NewInstance(NewCSSStyleDeclarationPrototype(), ctx)
}

func (p *CSSStyleDeclaration) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*CSSStyleDeclaration); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *CSSStyleDeclaration) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "display", "fontSize", "height", "width", "top", "bottom", "left", "right", "position": 
    return s, nil
  case "getPropertyValue":
    return values.NewFunction([]values.Value{s, s}, ctx), nil
  case "removeProperty":
    return values.NewMethodLikeFunction([]values.Value{s, s}, ctx), nil
  case "setProperty":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, nil},
      []values.Value{s, s, nil},
      []values.Value{s, s, s, nil},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *CSSStyleDeclaration) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  s := NewString(ctx)

  switch key {
  case "display", "height", "width", "top", "bottom", "left", "right", "position":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: CSSStyleDeclaration." + key + " not setable")
  }
}

func (p *CSSStyleDeclaration) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewUnconstructableClass(NewCSSStyleDeclarationPrototype(), ctx), nil
}
