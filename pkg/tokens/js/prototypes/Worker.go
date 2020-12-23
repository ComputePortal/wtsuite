package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Worker struct {
  BuiltinPrototype
}

func NewWorkerPrototype() values.Prototype {
  return &Worker{newBuiltinPrototype("Worker")}
}

func NewWorker(ctx context.Context) values.Value {
  return values.NewInstance(NewWorkerPrototype(), ctx)
}

func (p *Worker) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Worker); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Worker) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)

  switch key {
  case "onmessage", "onmessageerror":
    return nil, ctx.NewError("Error: only a setter")
  case "postMessage":
    return values.NewFunction([]values.Value{a, nil}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Worker) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  callback := values.NewFunction([]values.Value{NewMessageEvent(ctx), nil}, ctx)

  switch key {
  case "onmessage", "onmessageerror":
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: Worker." + key + " not setable")
  }
}

func (p *Worker) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  s := NewString(ctx)

  return values.NewClass([][]values.Value{
    []values.Value{s},
    []values.Value{s, s},
  }, NewWorkerPrototype(), ctx), nil
}
