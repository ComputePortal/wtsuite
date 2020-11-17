package prototypes

import (
  "../values"

  "../../context"
)

type Number struct {
  BuiltinPrototype
}

func NewNumberPrototype() values.Prototype {
  return &Number{newBuiltinPrototype("Number")}
}

func NewNumber(ctx context.Context) values.Value {
  return values.NewInstance(NewNumberPrototype(), ctx)
}

func IsNumber(v values.Value) bool {
  ctx := context.NewDummyContext()

  numberCheck := NewNumber(ctx)

  return numberCheck.Check(v, ctx) == nil
}

func (p *Number) IsUniversal() bool {
  return true
}

func (p *Number) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)
  i := NewInt(ctx)

  switch key {
  case "toExponential", "toFixed", "toPrecision", "toString":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{i, s},
    }, ctx), nil
  case "toLocaleString":
    opt := NewLocaleOptions(ctx)
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{s, s},
      []values.Value{s, opt, s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Number) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)
  b := NewBoolean(ctx)
  f := NewNumber(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "EPSILON", "MAX_VALUE", "MIN_VALUE", "NaN", "NEGATIVE_INFINITY", "POSITIVE_INFINITY":
    return f, nil
  case "MAX_SAFE_INTEGER", "MIN_SAFE_INTEGER":
    return i, nil
  case "isFinite", "isInteger", "isNaN", "isSafeInteger":
    return values.NewFunction([]values.Value{a, b}, ctx), nil
  case "parseFloat":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f},
      []values.Value{s, f},
    }, ctx), nil
  case "parseInt":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, i},
      []values.Value{s, i},
      []values.Value{f, i, i},
      []values.Value{s, i, i},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Number) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewClass([][]values.Value{
    []values.Value{NewNumber(ctx)},
    []values.Value{NewString(ctx)},
  }, NewNumberPrototype(), ctx), nil
}
