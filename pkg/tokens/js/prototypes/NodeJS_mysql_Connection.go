package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_mysql_Connection struct {
  BuiltinPrototype
}

func NewNodeJS_mysql_ConnectionPrototype() values.Prototype {
  return &NodeJS_mysql_Connection{newBuiltinPrototype("Connection")}
}

func NewNodeJS_mysql_Connection(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_mysql_ConnectionPrototype(), ctx)
}

func (p *NodeJS_mysql_Connection) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_mysql_Connection) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_mysql_Connection); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_mysql_Connection) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)
  ss := NewArray(s, ctx)

  switch key {
  case "connect", "release":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "end":
    return values.NewFunction([]values.Value{
      values.NewFunction([]values.Value{NewNodeJS_mysql_Error(ctx), nil}, ctx),
    }, ctx), nil
  case "query":
    callback := values.NewFunction([]values.Value{
      NewNodeJS_mysql_Error(ctx),
      values.NewAny(ctx),
      NewArray(NewNodeJS_mysql_FieldPacket(ctx), ctx),
      nil,
    }, ctx)

    q := NewNodeJS_mysql_Query(ctx)
    opt1 := NewConfigObject(map[string]values.Value{
      "sql": s,
      "timeout": i,
      "values": ss,
    }, ctx)

    opt2 := NewConfigObject(map[string]values.Value{
      "sql": s,
      "timeout": i,
    }, ctx)

    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{opt1, callback, q}, 
      []values.Value{opt2, ss, callback, q}, 
      []values.Value{s, callback, q}, 
      []values.Value{s, ss, callback, q}, 
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_mysql_Connection) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_mysql_ConnectionPrototype(), ctx), nil
}
