package prototypes

import (
	"../values"

	"../../context"
)

type Boolean struct {
  BuiltinPrototype
}

func NewBooleanPrototype() values.Prototype {
  return &Boolean{newBuiltinPrototype("Boolean")}
}

func NewBoolean(ctx context.Context) values.Value {
  return values.NewInstance(NewBooleanPrototype(), ctx)
}

func NewLiteralBoolean(v bool, ctx context.Context) values.Value {
  return values.NewLiteralBooleanInstance(NewBooleanPrototype(), v, ctx)
}

func IsBoolean(v values.Value) bool {
  ctx := context.NewDummyContext()
  
  booleanCheck := NewBoolean(ctx)

  return booleanCheck.Check(v, ctx) == nil
}

func (p *Boolean) IsUniversal() bool {
  return true
}

func (p *Boolean) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return nil, nil
}

func (p *Boolean) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()

  return values.NewClass(
    [][]values.Value{
      []values.Value{NewNumber(ctx)},
    }, NewBooleanPrototype(), ctx), nil
}
