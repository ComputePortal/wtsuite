package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type KeyboardEvent struct {
  targetless bool 

  AbstractEvent
}

func NewKeyboardEventPrototype(target values.Value) values.Prototype {
  return &KeyboardEvent{false, newAbstractEventPrototype("KeyboardEvent", target)}
}

func NewTargetlessKeyboardEventPrototype() values.Prototype {
  return &KeyboardEvent{true, newAbstractEventPrototype("KeyboardEvent", nil)}
}

func NewKeyboardEvent(target values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewKeyboardEventPrototype(target), ctx)
}

func (p *KeyboardEvent) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*KeyboardEvent); ok {
    return p.AbstractEvent.checkTarget(other.target, ctx)
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *KeyboardEvent) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  s := NewString(ctx)

  switch key {
  case "altKey", "ctrlKey", "metaKey", "shiftKey":
    return b, nil
  case "key":
    return s, nil
  default:
    if p.targetless && key == "target" {
      return nil, ctx.NewError("Error: targetless KeyboardEvent")
    } else {
      return p.AbstractEvent.GetInstanceMember(key, includePrivate, ctx)
    }
  }
}

func (p *KeyboardEvent) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  opt := NewConfigObject(map[string]values.Value{
    "altKey": b,
    "code": s,
    "ctrlKey": b,
    "isComposing": b,
    "key": s,
    "location": i,
    "metaKey": b,
    "repeat": b,
    "shiftKey": b,
  }, ctx)

  return values.NewClass([][]values.Value{
    []values.Value{s},
    []values.Value{s, opt},
  }, NewTargetlessKeyboardEventPrototype(), ctx), nil
}
