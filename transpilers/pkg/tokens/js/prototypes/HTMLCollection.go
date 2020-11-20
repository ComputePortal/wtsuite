package prototypes

import (
  "../values"

  "../../context"
)

type HTMLCollection struct {
  BuiltinPrototype
}

func NewHTMLCollectionPrototype() values.Prototype {
  return &HTMLCollection{newBuiltinPrototype("HTMLCollection")}
}

func NewHTMLCollection(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLCollectionPrototype(), ctx)
}

func (p *HTMLCollection) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLCollection); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLCollection) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)
  elem := NewHTMLElement(ctx)

  switch key {
  case ".getindex", "item":
    return values.NewFunction([]values.Value{i, elem}, ctx), nil
  case ".setindex":
    return values.NewFunction([]values.Value{i, elem, nil}, ctx), nil
  case "length":
    return i, nil
  case "nameItem":
    return values.NewFunction([]values.Value{s, elem}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *HTMLCollection) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLCollectionPrototype(), ctx), nil
}
