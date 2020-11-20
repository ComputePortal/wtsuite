package prototypes

import (
  "../values"

  "../../context"
)

type Element struct {
  BuiltinPrototype
}

func NewElementPrototype() values.Prototype {
  return &Element{newBuiltinPrototype("Element")}
}

func NewElement(ctx context.Context) values.Value {
  return values.NewInstance(NewElementPrototype(), ctx)
}

func (p *Element) GetParent() (values.Prototype, error) {
  return NewNodePrototype(), nil
}

func (p *Element) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Element); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Element) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  f := NewNumber(ctx)
  s := NewString(ctx)

  switch key {
  case "className", "id", "innerHTML", "tagName":
    return s, nil
  case "getAttribute":
    return values.NewFunction([]values.Value{s, s}, ctx), nil
  case "getBoundingClientRect":
    return values.NewFunction([]values.Value{NewDOMRect(ctx)}, ctx), nil
  case "hasAttribute": 
    return values.NewFunction([]values.Value{s, b}, ctx), nil
  case "removeAttribute":
    return values.NewFunction([]values.Value{s, nil}, ctx), nil
  case "scrollLeft", "scrollWidth", "scrollTop", "scrollHeight", "clientLeft", "clientWidth", "clientTop", "clientHeight":
    return f, nil
  case "scrollTo":
    return values.NewFunction([]values.Value{f, f, nil}, ctx), nil
  case "scrollIntoView":
    o := NewConfigObject(map[string]values.Value{
      "behavior": s,
      "block": s,
      "inline": s,
    }, ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{nil},
      []values.Value{b, nil},
      []values.Value{o, nil},
    }, ctx), nil
  case "setAttribute":
    return values.NewFunction([]values.Value{s, s, nil}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Element) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  f := NewNumber(ctx)
  s := NewString(ctx)

  switch key {
  case "className", "id", "innerHTML":
    return s.Check(arg, ctx)
  case "scrollLeft", "scrollTop":
    return f.Check(arg, ctx)
  default:
    return ctx.NewError("Error: Element." + key + " not setable")
  }
}

func (p *Element) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewElementPrototype(), ctx), nil
}
