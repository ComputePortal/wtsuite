package prototypes

import (
  "../values"

  "../../context"
)

type EventTarget struct {
  BuiltinPrototype
}

func NewEventTargetPrototype() values.Prototype {
  return &EventTarget{newBuiltinPrototype("EventTarget")}
}

func NewEventTarget(ctx context.Context) values.Value {
  return values.NewInstance(NewEventTargetPrototype(), ctx)
}

func (p *EventTarget) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*EventTarget); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *EventTarget) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)
  b := NewBoolean(ctx)
  s := NewString(ctx)

  switch key {
  case "addEventListener":
    // more specific events are possible
    return values.NewFunction([]values.Value{
      s, 
      values.NewFunction([]values.Value{NewEvent(a, ctx), nil}, ctx),
      nil,
    }, ctx), nil
  case "dispatchEvent":
    return values.NewMethodLikeFunction([]values.Value{NewEvent(a, ctx), b}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *EventTarget) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewEventTargetPrototype(), ctx), nil
}
