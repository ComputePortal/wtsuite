package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Document struct {
  BuiltinPrototype
}

func NewDocumentPrototype() values.Prototype {
  return &Document{newBuiltinPrototype("Document")}
}

func NewDocument(ctx context.Context) values.Value {
  return values.NewInstance(NewDocumentPrototype(), ctx)
}

func (p *Document) GetParent() (values.Prototype, error) {
  return NewEventTargetPrototype(), nil
}

func (p *Document) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Document); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Document) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  s := NewString(ctx)
  elem := NewHTMLElement(ctx)

  switch key {
  case "activeElement", "body", "documentElement":
    return elem, nil
  case "cookie", "referrer", "title":
    return s, nil
  case "createTextNode":
    return values.NewFunction([]values.Value{s, NewText(ctx)}, ctx), nil
  case "createElement":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewLiteralString("a", ctx), NewHTMLLinkElement(ctx)},
      []values.Value{NewLiteralString("canvas", ctx), NewHTMLCanvasElement(ctx)},
      []values.Value{NewLiteralString("img", ctx), NewHTMLImageElement(ctx)},
      []values.Value{NewLiteralString("input", ctx), NewHTMLInputElement(ctx)},
      []values.Value{s, elem}, // must come last, because first valid overload is used
    }, ctx), nil
  case "fonts":
    return NewFontFaceSet(ctx), nil
  case "getElementById", "querySelector":
    return values.NewFunction([]values.Value{s, elem}, ctx), nil
  case "hidden":
    return b, nil
  default:
    return nil, nil
  }
}

func (p *Document) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  s := NewString(ctx)

  switch key {
  case "cookie", "referrer", "title":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: document." + key + " not setable")
  }
}

func (p *Document) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewDocumentPrototype(), ctx), nil
}
