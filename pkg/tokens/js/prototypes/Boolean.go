package prototypes

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
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

func (p *Boolean) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Boolean); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
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
