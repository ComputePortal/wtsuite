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

func (p *MessageEvent) Name() string {
  return "MessageEvent"
}

func (p *MessageEvent) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*MessageEvent); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *MessageEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "data": 
    return NewObject(nil, ctx), nil
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
