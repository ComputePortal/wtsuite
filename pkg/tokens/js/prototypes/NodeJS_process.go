package prototypes

import (
  "../values"

  "../../context"
)

func FillNodeJS_processPackage(pkg values.Package) {
  ctx := context.NewDummyContext()
  i := NewInt(ctx)
  s := NewString(ctx)
  ss := NewArray(s, ctx)

  pkg.AddValue("argv", ss)
  pkg.AddValue("arg0", s)
  pkg.AddValue("exit", values.NewFunction([]values.Value{
    i, nil,
  }, ctx))
}
