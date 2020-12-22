package prototypes

import (
  "../values"

  "../../context"
)

type Blob struct {
  BuiltinPrototype
}

func NewBlobPrototype() values.Prototype {
  return &Blob{newBuiltinPrototype("Blob")}
}

func NewBlob(ctx context.Context) values.Value {
  return values.NewInstance(NewBlobPrototype(), ctx)
}

func IsBlob(v values.Value) bool {
  ctx := context.NewDummyContext()

  blobCheck := NewBlob(ctx)

  return blobCheck.Check(v, ctx) == nil
}

func (p *Blob) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Blob); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Blob) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)
  self := values.NewInstance(p, ctx)

  switch key {
  case "arrayBuffer":
    return values.NewFunction([]values.Value{
      NewPromise(NewArrayBuffer(ctx), ctx),
    }, ctx), nil
  case "size":
    return i, nil
  case "slice":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{self},
      []values.Value{i, self},
      []values.Value{i, i, self},
      []values.Value{i, i, s, self},
    }, ctx), nil
  case "type":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *Blob) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  s := NewString(ctx)
  o := NewConfigObject(map[string]values.Value{
    "type": s,
  }, ctx)

  return values.NewClass([][]values.Value{
    []values.Value{s},
    []values.Value{s, o},
    []values.Value{NewArray(s, ctx)}, 
    []values.Value{NewArray(s, ctx), o}, 
    []values.Value{NewArray(NewArrayBuffer(ctx), ctx)}, 
    []values.Value{NewArray(NewArrayBuffer(ctx), ctx), o}, 
    []values.Value{NewArray(NewBlob(ctx), ctx)}, 
    []values.Value{NewArray(NewBlob(ctx), ctx), o}, 
  }, NewBlobPrototype(), ctx), nil
}
