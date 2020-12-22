package prototypes

import (
  "../values"

  "../../context"
)

func FillMathPackage(pkg values.Package) {
  ctx := context.NewDummyContext()
  f := NewNumber(ctx)
  i := NewInt(ctx)
  one2one := values.NewFunction([]values.Value{f, f}, ctx)
  round := values.NewFunction([]values.Value{f, i}, ctx)
  minmax := values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, i, i},
      []values.Value{i, i, i, i},
      []values.Value{i, i, i, i, i},
      []values.Value{i, i, i, i, i, i},
      []values.Value{i, i, i, i, i, i, i},
      []values.Value{f, f, f},
      []values.Value{f, f, f, f},
      []values.Value{f, f, f, f, f},
      []values.Value{f, f, f, f, f, f},
      []values.Value{f, f, f, f, f, f, f}, // should be enough
    }, ctx)

  pkg.AddValue("E", f)
  pkg.AddValue("LN2", f)
  pkg.AddValue("LN10", f)
  pkg.AddValue("LOG2E", f)
  pkg.AddValue("LOG10E", f)
  pkg.AddValue("PI", f)
  pkg.AddValue("SQRT1_2", f)
  pkg.AddValue("SQRT2", f)

  pkg.AddValue("abs", one2one)
  pkg.AddValue("acos", one2one)
  pkg.AddValue("acosh", one2one)
  pkg.AddValue("asin", one2one)
  pkg.AddValue("asinh", one2one)
  pkg.AddValue("atan", one2one)
  pkg.AddValue("atanh", one2one)
  pkg.AddValue("cbrt", one2one)
  pkg.AddValue("cos", one2one)
  pkg.AddValue("cosh", one2one)
  pkg.AddValue("exp", one2one)
  pkg.AddValue("expm1", one2one)
  pkg.AddValue("fround", one2one)
  pkg.AddValue("log", one2one)
  pkg.AddValue("log10", one2one)
  pkg.AddValue("log1p", one2one)
  pkg.AddValue("log2", one2one)
  pkg.AddValue("sin", one2one)
  pkg.AddValue("sinh", one2one)
  pkg.AddValue("sqrt", one2one)
  pkg.AddValue("tan", one2one)
  pkg.AddValue("tanh", one2one)

  pkg.AddValue("atan2", values.NewFunction([]values.Value{f, f, f}, ctx))

  pkg.AddValue("ceil", round)
  pkg.AddValue("floor", round)
  pkg.AddValue("round", round)
  pkg.AddValue("sign", round)
  pkg.AddValue("trunc", round)

  pkg.AddValue("hypot", values.NewOverloadedFunction([][]values.Value{
      []values.Value{f},
      []values.Value{f, f},
      []values.Value{f, f, f},
      []values.Value{f, f, f, f},
      []values.Value{f, f, f, f, f}, // should be enough
    }, ctx))

  pkg.AddValue("min", minmax)
  pkg.AddValue("max", minmax)

  pkg.AddValue("pow", values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f, f},
      []values.Value{i, i, i},
    }, ctx))

  pkg.AddValue("random", values.NewFunction([]values.Value{f}, ctx))
}
