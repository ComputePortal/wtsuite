package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_stream_Readable struct {
  BuiltinPrototype
}

func NewNodeJS_stream_ReadablePrototype() values.Prototype {
  return &NodeJS_stream_Readable{newBuiltinPrototype("Readable")}
}

func NewNodeJS_stream_Readable(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_stream_ReadablePrototype(), ctx)
}

func (p *NodeJS_stream_Readable) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_stream_Readable) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_stream_Readable); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_stream_Readable) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "addListener":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewLiteralString("data", ctx), values.NewFunction([]values.Value{
        s, nil,
      }, ctx), nil},
      []values.Value{s, values.NewFunction([]values.Value{nil}, ctx), nil},
    }, ctx), nil
  case "read":
    i := NewInt(ctx)
    buf := NewNodeJS_Buffer(ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{buf},
      []values.Value{i, buf},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_stream_Readable) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_stream_ReadablePrototype(), ctx), nil
}
