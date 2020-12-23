package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Location struct {
  BuiltinPrototype
}

func NewLocationPrototype() values.Prototype {
  return &Location{newBuiltinPrototype("Location")}
}

func NewLocation(ctx context.Context) values.Value {
  return values.NewInstance(NewLocationPrototype(), ctx)
}

func (p *Location) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Location); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Location) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "hash", "href", "origin", "pathname", "search": 
    return s, nil
  default:
    return nil, nil
  }
}

func (p *Location) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  s := NewString(ctx)

  switch key {
  case "hash", "href":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: Location." + key + " not setable")
  }
}

func (p *Location) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewUnconstructableClass(NewLocationPrototype(), ctx), nil
}
