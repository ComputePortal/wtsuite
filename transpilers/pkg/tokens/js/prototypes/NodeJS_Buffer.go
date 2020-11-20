package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_Buffer struct {
  BuiltinPrototype
}

func NewNodeJS_BufferPrototype() values.Prototype {
  return &NodeJS_Buffer{newBuiltinPrototype("Buffer")}
}

func NewNodeJS_Buffer(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_BufferPrototype(), ctx)
}

func (p *NodeJS_Buffer) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_Buffer); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_Buffer) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewString(ctx)
  s := NewString(ctx)

  switch key {
  case "toString":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{s, s},
      []values.Value{s, i, s},
      []values.Value{s, i, i, s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_Buffer) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)
  self := NewNodeJS_Buffer(ctx)

  switch key {
  case "concat":
    return values.NewFunction([]values.Value{NewArray(self, ctx), self}, ctx), nil
  case "from":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{self, self},
      []values.Value{NewArray(NewInt(ctx), ctx), self},
      []values.Value{s, self},
      []values.Value{s, s, self},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_Buffer) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()
  return values.NewUnconstructableClass(NewNodeJS_BufferPrototype(), ctx), nil
}
