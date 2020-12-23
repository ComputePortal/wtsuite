package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type DataView struct {
  BuiltinPrototype
}

func NewDataViewPrototype() values.Prototype {
  return &DataView{newBuiltinPrototype("DataView")}
}

func NewDataView(ctx context.Context) values.Value {
  return values.NewInstance(NewDataViewPrototype(), ctx)
}

func (p *DataView) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*DataView); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *DataView) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  f := NewNumber(ctx)
  big := NewBigInt(ctx)

  switch key {
  case "getInt8", "getUint8", "getInt16", "getUint16", "getInt32", "getUint32":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, i},
      []values.Value{i, b, i},
    }, ctx), nil
  case "getBigInt64", "getBigUint64":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, big},
      []values.Value{i, b, big},
    }, ctx), nil
  case "getFloat32", "getFloat64":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, f},
      []values.Value{i, b, f},
    }, ctx), nil
  case "setInt8", "setUint8", "setInt16", "setUint16", "setInt32", "setUint32":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, i, nil},
      []values.Value{i, i, b, nil},
    }, ctx), nil
  case "setBigInt64", "setBigUint64":

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, big, nil},
      []values.Value{i, big, b, nil},
      []values.Value{i, i, nil},
      []values.Value{i, i, b, nil},
    }, ctx), nil
  case "setFloat32", "setFloat64":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, f, nil},
      []values.Value{i, f, b, nil},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *DataView) GetClassValue() (*values.Class, error) {
  ctx := context.NewDummyContext()

  return values.NewClass([][]values.Value{
    []values.Value{NewArrayBuffer(ctx)},
  }, NewDataViewPrototype(), ctx), nil
}
