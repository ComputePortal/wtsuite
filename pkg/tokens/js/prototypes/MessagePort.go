package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type MessagePort struct {
  BuiltinPrototype
}

func NewMessagePortPrototype() values.Prototype {
  return &MessagePort{newBuiltinPrototype("MessagePort")}
}

func NewMessagePort(ctx context.Context) values.Value {
  return values.NewInstance(NewMessagePortPrototype(), ctx)
}

func (p *MessagePort) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*MessagePort); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *MessagePort) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)

  switch key {
  case "onmessage", "onmessageerror":
    return nil, ctx.NewError("Error: is only a setter")
  case "close", "start":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "postMessage":
    return values.NewFunction([]values.Value{a, nil}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *MessagePort) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  callback := values.NewFunction([]values.Value{NewMessageEvent(ctx), nil}, ctx)

  switch key {
  case "onmessage", "onmessageerror":
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: MessagePort." + key + " not setable")
  }
}

func (p *MessagePort) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewUnconstructableClass(NewMessagePortPrototype(), ctx), nil
}
