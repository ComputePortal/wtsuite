package prototypes

import (
  "../values"

  "../../context"
)

type FileReader struct {
  BuiltinPrototype
}

func NewFileReaderPrototype() values.Prototype {
  return &FileReader{newBuiltinPrototype("FileReader")}
}

func NewFileReader(ctx context.Context) values.Value {
  return values.NewInstance(NewFileReaderPrototype(), ctx)
}

func (p *FileReader) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "readAsArrayBuffer":
    return values.NewFunction([]values.Value{NewBlob(ctx)}, ctx), nil
  case "result":
    return NewArrayBuffer(ctx), nil
  default:
    return nil, nil
  }
}

func (p *FileReader) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  switch key {
  case "onload", "onerror":
    callback := values.NewFunction([]values.Value{NewEvent(nil, ctx), nil}, ctx)
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: FileReader." + key + " not setable")
  }
}

func (p *FileReader) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewClass([][]values.Value{
    []values.Value{},
  }, NewFileReaderPrototype(), ctx), nil
}
