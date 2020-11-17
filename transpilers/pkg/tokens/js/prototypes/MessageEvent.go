package prototypes

import (
  "../values"

  "../../context"
)

type MessageEvent struct {
  AbstractEvent
}

func NewMessageEventPrototype() values.Prototype {
  ctx := context.NewDummyContext()
  return &MessageEvent{newAbstractEventPrototype("MessageEvent", NewMessagePort(ctx))}
}

func NewMessageEvent(ctx context.Context) values.Value {
  return values.NewInstance(NewMessageEventPrototype(), ctx)
}

func (p *MessageEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)

  switch key {
  case "data": 
    return a, nil
  case "ports":
    return NewArray(p.target, ctx), nil
  default:
    return p.AbstractEvent.GetInstanceMember(key, includePrivate, ctx)
  }
}

func (p *MessageEvent) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewMessageEventPrototype(), ctx), nil
}
