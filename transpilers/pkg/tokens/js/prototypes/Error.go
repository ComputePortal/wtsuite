package prototypes

import (
  "../values"

  "../../context"
)

type Error struct {
  BuiltinPrototype
}

func NewErrorPrototype() values.Prototype {
  return &Error{newBuiltinPrototype("Error")}
}

func NewError(ctx context.Context) values.Value {
  return values.NewInstance(NewErrorPrototype(), ctx)
}

func IsError(v values.Value) bool {
  ctx := context.NewDummyContext()

  errorCheck := NewError(ctx)

  return errorCheck.Check(v, ctx) == nil
}

func (p *Error) IsUniversal() bool {
  return true
}

func (p *Error) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "message":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *Error) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()

  return values.NewClass([][]values.Value{
    []values.Value{},
    []values.Value{NewString(ctx)},
  }, NewErrorPrototype(), ctx), nil
}
