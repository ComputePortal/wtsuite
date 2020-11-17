package prototypes

import (
  "../values"

  "../../context"
)

type WheelEvent struct {
  AbstractEvent
}

func NewWheelEventPrototype(target values.Value) values.Prototype {
  return &WheelEvent{newAbstractEventPrototype("WheelEvent", target)}
}

func NewWheelEvent(target values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewWheelEventPrototype(target), ctx)
}

func (p *WheelEvent) GetParent() (values.Prototype, error) {
  return NewMouseEventPrototype(p.target), nil
}

func (p *WheelEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  f := NewNumber(ctx)

  switch key {
  case "deltaX", "deltaY", "deltaZ":
    return f, nil
  default:
    return nil, nil
  }
}

func (p *WheelEvent) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewUnconstructableClass(NewWheelEventPrototype(values.NewAll(ctx)), ctx), nil
}
