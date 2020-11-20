package prototypes

import (
	"../values"

	"../../context"
)

type Int struct {
  BuiltinPrototype
}

func NewIntPrototype() values.Prototype {
  return &Int{newBuiltinPrototype("Int")}
}

func NewInt(ctx context.Context) values.Value {
  return values.NewInstance(NewIntPrototype(), ctx)
}

func NewLiteralInt(v int, ctx context.Context) values.Value {
  return values.NewLiteralIntInstance(NewIntPrototype(), v, ctx)
}

func IsInt(v values.Value) bool {
  ctx := context.NewDummyContext()

  intCheck := NewInt(ctx)

  return intCheck.Check(v, ctx) == nil
}

func (p *Int) GetParent() (values.Prototype, error) {
  return NewNumberPrototype(), nil
}

func (p *Int) IsUniversal() bool {
  return true
}

func (p *Int) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Int); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Int) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "toString":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{i, s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Int) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()
  f := NewNumber(ctx)
  s := NewString(ctx)

  return values.NewClass(
    [][]values.Value{
      []values.Value{f},
      []values.Value{s},
    }, NewIntPrototype(), ctx), nil
}
