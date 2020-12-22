package prototypes

import (
  "../values"

  "../../context"
)

func FillNodeJS_httpPackage(pkg values.Package) {
  pkg.AddPrototype(NewNodeJS_http_IncomingMessagePrototype())
  pkg.AddPrototype(NewNodeJS_http_ServerPrototype())
  pkg.AddPrototype(NewNodeJS_http_ServerResponsePrototype())

  ctx := context.NewDummyContext()

  pkg.AddValue("createServer", values.NewFunction([]values.Value{
    NewNodeJS_http_Server(ctx),
  }, ctx))
}
