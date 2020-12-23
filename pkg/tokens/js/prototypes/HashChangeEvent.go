package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type HashChangeEvent struct {
  AbstractEvent
}

func NewHashChangeEventPrototype() values.Prototype {
  ctx := context.NewDummyContext()
  return &HashChangeEvent{newAbstractEventPrototype("HashChangeEvent", NewWindow(ctx))}
}

func NewHashChangeEvent(ctx context.Context) values.Value {
  return values.NewInstance(NewHashChangeEventPrototype(), ctx)
}

func (p *HashChangeEvent) Name() string {
  return "HashChangeEvent"
}

func (p *HashChangeEvent) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HashChangeEvent); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HashChangeEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "newURL", "oldURL":
    return s, nil
  default:
    return p.AbstractEvent.GetInstanceMember(key, includePrivate, ctx)
  }
}

func (p *HashChangeEvent) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHashChangeEventPrototype(), ctx), nil
}
