package prototypes

import (
  "../values"

  "../../context"
)

type TypedArray interface {
  values.Prototype

  isUnsigned() bool

  numBits() int

  getContent() values.Value
}

type AbstractTypedArray struct {
  unsigned bool

  nBits int

  content values.Value

  BuiltinPrototype
}

func (p *AbstractTypedArray) isUnsigned() bool {
  return p.unsigned
}

func (p *AbstractTypedArray) numBits() int {
  return p.nBits
}

func (p *AbstractTypedArray) getContent() values.Value {
  return p.content
}

func newAbstractTypedArrayPrototype(name string, unsigned bool, nBits int, content values.Value) AbstractTypedArray {
  if content == nil {
    panic("content can't be nil (unlike Array)")
  }

  return AbstractTypedArray{unsigned, nBits, content, newBuiltinPrototype(name)}
}

func NewTypedArray(ctx context.Context) values.Value {
  proto := newAbstractTypedArrayPrototype("TypedArray", false, 0, NewNumber(ctx))
  return values.NewInstance(&proto, ctx)
}

func IsTypedArray(v values.Value) bool {
  ctx := context.NewDummyContext()
  
  typedArrayCheck := NewTypedArray(ctx)

  return typedArrayCheck.Check(v, ctx) == nil
}

func (p *AbstractTypedArray) GetParent() (values.Prototype, error) {
  return NewArrayPrototype(p.content), nil
}

func (p *AbstractTypedArray) IsUniversal() bool {
  return true
}

func CheckTypedArray(p TypedArray, other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(TypedArray); ok {
    thisContent := p.getContent()
    otherContent := other.getContent()

    if p.numBits() == 0 {
      return nil
    }

    if thisContent.Check(otherContent, ctx) != nil || p.numBits() != other.numBits() || p.isUnsigned() != other.isUnsigned() {
      return ctx.NewError("Error: expected " + p.Name() + ", got " + other_.Name())
    }

    return nil
  } else {
    return ctx.NewError("Error: expected TypedArray, got " + other_.Name())
  }
}

func GetTypedArrayInstanceMember(p TypedArray, key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  arr := NewArray(p.getContent(), ctx)
  self := values.NewInstance(p, ctx)

  switch key {
  case "BYTES_PER_ELEMENT":
    return i, nil
  case "buffer":
    return NewArrayBuffer(ctx), nil
  case "set":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{arr, nil},
      []values.Value{arr, i, nil},
    }, ctx), nil
  case "slice":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{self},
      []values.Value{i, self},
      []values.Value{i, i, self},
    }, ctx), nil
  case "subarray":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{self},
      []values.Value{i, self},
      []values.Value{i, i, self},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func GetTypedArrayClassMember(p TypedArray, key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  self := values.NewInstance(p, ctx)
  content := p.getContent()

  switch key {
  case "from":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewSet(content, ctx), self},
      []values.Value{NewArray(content, ctx), self},
    }, ctx), nil
  default: 
    return nil, nil
  }
}

func GetTypedArrayClassValue(p TypedArray) (*values.Class, error) {
  ctx := p.Context()

  i := NewInt(ctx)
  content := p.getContent()

  return values.NewClass([][]values.Value{
    []values.Value{i},
    []values.Value{NewArray(content, ctx)},
    []values.Value{NewArrayBuffer(ctx)},
  }, p, ctx), nil
}

