package prototypes

import (
  "../values"

  "../../context"
)

func FillNodeJS_nodemailerPackage(pkg values.Package) {
  pkg.AddPrototype(NewNodeJS_nodemailer_SMTPTransportPrototype())

  ctx := context.NewDummyContext()
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  opt := NewConfigObject(map[string]values.Value{
    "host": s,
    "port": i,
    "secure": b,
    "auth": NewObject(map[string]values.Value{
      "user": s,
      "pass": s,
    }, ctx),
  }, ctx)

  pkg.AddValue("createTransport", values.NewFunction([]values.Value{
    opt, NewNodeJS_nodemailer_SMTPTransport(ctx),
  }, ctx))
}
