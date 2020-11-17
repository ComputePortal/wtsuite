package prototypes

import (
  "../values"

  "../../context"
)

type Math struct {
  BuiltinPrototype
}

func NewMathPrototype() values.Prototype {
  return &Math{newBuiltinPrototype("Math")}
}

func NewMath(ctx context.Context) values.Value {
  return values.NewInstance(NewMathPrototype(), ctx)
}

func (p *Math) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  f := NewNumber(ctx)
  i := NewInt(ctx)

  switch key {
  case "E", "LN2", "LN10", "LOG2E", "LOG10E", "PI", "SQRT1_2", "SQRT2":
    return f, nil
  case "abs", "acos", "acosh", "asin", "asinh", "atan", "atanh", "cbrt", "cos", "cosh", "exp", "expm1", "fround", "log", "log10", "log1p", "log2", "sin", "sinh", "sqrt", "tan", "tanh":
    return values.NewFunction([]values.Value{f, f}, ctx), nil
  case "atan2":
    return values.NewFunction([]values.Value{f, f, f}, ctx), nil
  case "ceil", "floor", "round", "sign", "trunc":
    return values.NewFunction([]values.Value{f, i}, ctx), nil
  case "hypot":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f},
      []values.Value{f, f},
      []values.Value{f, f, f},
      []values.Value{f, f, f, f},
      []values.Value{f, f, f, f, f}, // should be enough
    }, ctx), nil
  case "min", "max":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f},
      []values.Value{f, f, f},
      []values.Value{f, f, f, f},
      []values.Value{f, f, f, f, f},
      []values.Value{f, f, f, f, f, f}, // should be enough
    }, ctx), nil
  case "pow":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f, f},
      []values.Value{i, i, i},
    }, ctx), nil
  case "random":
    return values.NewFunction([]values.Value{f}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Math) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewMathPrototype(), ctx), nil
}
