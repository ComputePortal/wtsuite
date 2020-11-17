package prototypes

import (
  "../values"

  "../../context"
)

func FillNodeJS_fsPackage(pkg values.Package) {
  ctx := context.NewDummyContext()
  b := NewBoolean(ctx)
  s := NewString(ctx)
  buf := NewNodeJS_Buffer(ctx)

  pkg.AddValue("existsSync", values.NewFunction([]values.Value{s, b}, ctx))

  pkg.AddValue("readFileSync", values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, buf},
      []values.Value{s, NewLiteralString("buffer", ctx), buf},
      []values.Value{s, s, s},
    }, ctx))

  pkg.AddValue("unlinkSync", values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, nil},
      []values.Value{buf, nil},
    }, ctx))
}
