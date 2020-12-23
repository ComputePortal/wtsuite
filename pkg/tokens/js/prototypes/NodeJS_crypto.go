package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

func FillNodeJS_cryptoPackage(pkg values.Package) {
  pkg.AddPrototype(NewNodeJS_crypto_CipherPrototype())
  pkg.AddPrototype(NewNodeJS_crypto_DecipherPrototype())

  ctx := context.NewDummyContext()
  i := NewInt(ctx)
  s := NewString(ctx)
  buf := NewNodeJS_Buffer(ctx)

  pkg.AddValue("createCipheriv", values.NewFunction([]values.Value{
    s, buf, buf, NewNodeJS_crypto_Cipher(ctx),
  }, ctx))

  pkg.AddValue("createDecipheriv", values.NewFunction([]values.Value{
    s, buf, buf, NewNodeJS_crypto_Decipher(ctx),
  }, ctx))

  pkg.AddValue("randomBytes", values.NewFunction([]values.Value{ 
    i, buf,
  }, ctx))
}
