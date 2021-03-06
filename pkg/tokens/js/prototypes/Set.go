package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Set struct {
  content values.Value // if nil, then any

  BuiltinPrototype
}

func NewSetPrototype(content values.Value) values.Prototype {
  return &Set{content, newBuiltinPrototype("Set")}
}

func NewSet(content values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewSetPrototype(content), ctx)
}

// what if other inherits from set?
func (p *Set) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Set); ok {
    if p.content == nil {
      return nil
    } else if other.content == nil {
      return ctx.NewError("Error: expected Set<" + p.content.TypeName() + ">, got Set<any>")
    } else if p.content.Check(other.content, ctx) != nil {
      return ctx.NewError("Error: expected Set<" + p.content.TypeName() + ">, got Set<" + other.content.TypeName() + ">")
    } else {
      return nil
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Set) getContentValue() values.Value {
  if p.content == nil {
    return values.NewAny(context.NewDummyContext())
  } else {
    return p.content
  }
}

func (p *Set) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  content := values.NewContextValue(p.getContentValue(), ctx)
  self := values.NewInstance(p, ctx)

  switch key {
  case ".content", ".getof":
    return content, nil
  case "add":
    return values.NewMethodLikeFunction([]values.Value{content, self}, ctx), nil
  case "clear":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "delete":
    return values.NewMethodLikeFunction([]values.Value{content, b}, ctx), nil
  case "has":
    return values.NewFunction([]values.Value{content, b}, ctx), nil
  case "size":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *Set) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewClass(
    [][]values.Value{
      []values.Value{},
    }, NewSetPrototype(values.NewAny(ctx)), ctx), nil
}
