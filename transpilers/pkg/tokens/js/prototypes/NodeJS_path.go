package prototypes

import (
  "../values"

  "../../context"
)

func FillNodeJS_pathPackage(pkg values.Package) {
  ctx := context.NewDummyContext()
  s := NewString(ctx)

  pkg.AddValue("join", values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, s, s},
      []values.Value{s, s, s, s},
      []values.Value{s, s, s, s, s},
      []values.Value{s, s, s, s, s, s},
      []values.Value{s, s, s, s, s, s, s}, // should be enough
    }, ctx))
}
