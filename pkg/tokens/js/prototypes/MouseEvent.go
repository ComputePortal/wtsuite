package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type MouseEvent struct {
  AbstractEvent
}

func NewMouseEventPrototype() values.Prototype {
  ctx := context.NewDummyContext()
  return &MouseEvent{newAbstractEventPrototype("MouseEvent", NewHTMLElement(ctx))}
}

func NewMouseEvent(ctx context.Context) values.Value {
  return values.NewInstance(NewMouseEventPrototype(), ctx)
}

func (p *MouseEvent) Name() string {
  return "MouseEvent"
}

func (p *MouseEvent) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*MouseEvent); ok {
    return p.AbstractEvent.checkTarget(other.target, ctx)
  } else {
    return checkParent(p, other_, ctx)
  }
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
  return values.NewUnconstructableClass(NewMouseEventPrototype(), ctx), nil
}
