package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

func FillJSONPackage(pkg values.Package) {
  ctx := context.NewDummyContext()
  s := NewString(ctx)
  o := NewObject(nil, ctx)

  pkg.AddValue("stringify", values.NewFunction([]values.Value{o, s}, ctx))
  pkg.AddValue("parse", values.NewFunction([]values.Value{s, o}, ctx))
}
