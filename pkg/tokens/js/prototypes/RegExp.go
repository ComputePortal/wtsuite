package prototypes

import (
  "../values"

  "../../context"
)

type RegExp struct {
  BuiltinPrototype
}

func NewRegExpPrototype() values.Prototype {
  return &RegExp{newBuiltinPrototype("RegExp")}
}

func NewRegExp(ctx context.Context) values.Value {
  return values.NewInstance(NewRegExpPrototype(), ctx)
}

func IsRegExp(v values.Value) bool {
  ctx := context.NewDummyContext()

  regexpCheck := NewRegExp(ctx)

  return regexpCheck.Check(v, ctx) == nil
}

func (p *RegExp) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*RegExp); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *RegExp) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "exec":
    return values.NewFunction([]values.Value{s, NewRegExpArray(ctx)}, ctx), nil
  case "global", "ignoreCase", "multiline":
    return b, nil
  case "lastIndex":
    return i, nil
  case "source":
    return s, nil
  case "test":
    return values.NewFunction([]values.Value{s, b}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *RegExp) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  s := NewString(ctx)

  return values.NewClass([][]values.Value{
    []values.Value{s},
    []values.Value{s, s},
  }, NewRegExpPrototype(), ctx), nil
}
