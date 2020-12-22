package prototypes

import (
  "../values"
)

func FillNodeJS_streamPackage(pkg values.Package) {
  pkg.AddPrototype(NewNodeJS_stream_ReadablePrototype())
}
