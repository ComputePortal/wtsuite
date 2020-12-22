package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_EventEmitter struct {
  BuiltinPrototype
}

func NewNodeJS_EventEmitterPrototype() values.Prototype {
  return &NodeJS_EventEmitter{newBuiltinPrototype("EventEmitter")}
}

func NewNodeJS_EventEmitter(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_EventEmitterPrototype(), ctx)
}

func (p *NodeJS_EventEmitter) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_EventEmitter); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_EventEmitter) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "addListener":
    return values.NewFunction([]values.Value{
      s,
      values.NewFunction([]values.Value{nil}, ctx),
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_EventEmitter) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()
  return values.NewUnconstructableClass(NewNodeJS_EventEmitterPrototype(), ctx), nil
}
