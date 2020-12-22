package prototypes

import (
  "strings"

  "../values"

  "../../context"
)

type AbstractEvent struct {
  target values.Value // if nil, then any target

  BuiltinPrototype
}

type Event struct {
  AbstractEvent
}

func newAbstractEventPrototype(name string, target values.Value) AbstractEvent {
  return AbstractEvent{target, newBuiltinPrototype(name)}
}

func NewEventPrototype(target values.Value) values.Prototype {
  return &Event{newAbstractEventPrototype("Event", target)}
}

func NewEvent(target values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewEventPrototype(target), ctx)
}

func (p *AbstractEvent) Name() string {
  var b strings.Builder

  b.WriteString(p.name)

  if p.target != nil {
    b.WriteString("<")
    b.WriteString(p.target.TypeName())
    b.WriteString(">")
  }

  return b.String()
}

func (p *Event) GetParent() (values.Prototype, error) {
  return nil, nil
}

func (p *AbstractEvent) GetParent() (values.Prototype, error) {
  return NewEventPrototype(p.target), nil
}

func (p *Event) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Event); ok {
    return p.AbstractEvent.checkTarget(other.target, ctx)
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *AbstractEvent) checkTarget(otherTarget values.Value, ctx context.Context) error {
  if p.target == nil {
    return nil
  } else if otherTarget == nil && !values.IsAny(p.target) {
    return ctx.NewError("Error: expected " + p.Name() + ", got Event<any>")
  } else if p.target.Check(otherTarget, ctx) != nil {
    return ctx.NewError("Error: expected " + p.Name() + ", got Event<" + otherTarget.TypeName() + ">")
  } else {
    return nil
  }
}

func (p *AbstractEvent) getTargetValue() values.Value {
  if p.target == nil {
    return values.NewAny(p.Context())
  } else {
    return p.target
  }
}

func (p *AbstractEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "target":
    target := p.getTargetValue()
    return values.NewContextValue(target, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Event) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "preventDefault", "stopPropagation", "stopImmediatePropagation":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  default:
    return p.AbstractEvent.GetInstanceMember(key, includePrivate, ctx)
  }
}

func (p *Event) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  a := values.NewAny(ctx)
  b := NewBoolean(ctx)
  s := NewString(ctx)
  o := NewConfigObject(map[string]values.Value{
    "bubbles": b,
    "cancelable": b,
    "composed": b,
  }, ctx)

  return values.NewClass(
    [][]values.Value{
      []values.Value{s},
      []values.Value{s, o},
    }, NewEventPrototype(a), ctx), nil
}
