package prototypes

import (
  "../values"

  "../../context"
)

type SharedWorker struct {
  BuiltinPrototype
}

func NewSharedWorkerPrototype() values.Prototype {
  return &SharedWorker{newBuiltinPrototype("SharedWorker")}
}

func NewSharedWorker(ctx context.Context) values.Value {
  return values.NewInstance(NewSharedWorkerPrototype(), ctx)
}

func IsSharedWorker(v values.Value) bool {
  ctx := context.NewDummyContext()

  checkVal := NewSharedWorker(ctx)

  return checkVal.Check(v, ctx) == nil
}

func (p *SharedWorker) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "port":
    return NewMessagePort(ctx), nil
  default:
    return nil, nil
  }
}

func (p *SharedWorker) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  s := NewString(ctx)

  return values.NewClass([][]values.Value{
    []values.Value{s},
    []values.Value{s, s},
  }, NewSharedWorkerPrototype(), ctx), nil
}
