package prototypes

import (
  "../values"

  "../../context"
)

func FillNodeJS_mysqlPackage(pkg values.Package) {
  pkg.AddPrototype(NewNodeJS_mysql_ConnectionPrototype())
  pkg.AddPrototype(NewNodeJS_mysql_ErrorPrototype())
  pkg.AddPrototype(NewNodeJS_mysql_FieldPacketPrototype())
  pkg.AddPrototype(NewNodeJS_mysql_PoolPrototype())
  pkg.AddPrototype(NewNodeJS_mysql_QueryPrototype())

  ctx := context.NewDummyContext()

  i := NewInt(ctx)
  s := NewString(ctx)

  connOpt := NewConfigObject(map[string]values.Value{
    "host": s,
    "post": i,
    "user": s,
    "password": s,
    "database": s,
  }, ctx)

  pkg.AddValue("createConnection", values.NewFunction([]values.Value{
    connOpt, NewNodeJS_mysql_Connection(ctx),
  }, ctx))

  poolOpt := NewConfigObject(map[string]values.Value{
    "connectionLimit": i,
    "host": s,
    "user": s,
    "password": s,
    "database": s,
  }, ctx)

  pkg.AddValue("createPool", values.NewFunction([]values.Value{
    poolOpt, NewNodeJS_mysql_Pool(ctx),
  }, ctx))
}
