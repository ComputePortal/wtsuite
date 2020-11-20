package prototypes

import (
  "../values"

  "../../context"
)

type IDBVersionChangeEvent struct {
  AbstractEvent
}

func NewIDBVersionChangeEventPrototype() values.Prototype {
  ctx := context.NewDummyContext()
  return &IDBVersionChangeEvent{newAbstractEventPrototype("IDBVersionChangeEvent", NewIDBOpenDBRequest(ctx))}
}

func NewIDBVersionChangeEvent(ctx context.Context) values.Value {
  return values.NewInstance(NewIDBVersionChangeEventPrototype(), ctx)
}

func (p *IDBVersionChangeEvent) Name() string {
  return "IDBVersionChangeEvent"
}

func (p *IDBVersionChangeEvent) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*IDBVersionChangeEvent); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBVersionChangeEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)

  switch key {
  case "newVersion", "oldVersion":
    return i, nil
  default:
    return p.AbstractEvent.GetInstanceMember(key, includePrivate, ctx)
  }
}

func (p *IDBVersionChangeEvent) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBVersionChangeEventPrototype(), ctx), nil
}
