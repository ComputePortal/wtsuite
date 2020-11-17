package prototypes

import (
  "../values"

  "../../context"
)

type MouseEvent struct {
  AbstractEvent
}

func NewMouseEventPrototype(target values.Value) values.Prototype {
  return &MouseEvent{newAbstractEventPrototype("MouseEvent", target)}
}

func NewMouseEvent(target values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewMouseEventPrototype(target), ctx)
}

func (p *MouseEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  f := NewNumber(ctx)

  switch key {
  case "altKey", "ctrlKey", "metaKey", "shiftKey":
    return b, nil
  case "clientX", "clientY":
    return f, nil
  default:
    return nil, nil
  }
}

func (p *MouseEvent) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()
  return values.NewUnconstructableClass(NewMouseEventPrototype(values.NewAll(ctx)), ctx), nil
}
