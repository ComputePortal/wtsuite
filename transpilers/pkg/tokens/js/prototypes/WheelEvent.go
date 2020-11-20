package prototypes

import (
  "../values"

  "../../context"
)

type WheelEvent struct {
  AbstractEvent
}

func NewWheelEventPrototype() values.Prototype {
  ctx := context.NewDummyContext()
  return &WheelEvent{newAbstractEventPrototype("WheelEvent", NewHTMLElement(ctx))}
}

func NewWheelEvent(ctx context.Context) values.Value {
  return values.NewInstance(NewWheelEventPrototype(), ctx)
}

func (p *WheelEvent) Name() string {
  return "WheelEvent"
}

func (p *WheelEvent) GetParent() (values.Prototype, error) {
  return NewMouseEventPrototype(), nil
}

func (p *WheelEvent) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*WheelEvent); ok {
    return p.AbstractEvent.checkTarget(other.target, ctx)
  } else {
    return checkParent(p, other_, ctx)
  }
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

  return values.NewUnconstructableClass(NewWheelEventPrototype(), ctx), nil
}
