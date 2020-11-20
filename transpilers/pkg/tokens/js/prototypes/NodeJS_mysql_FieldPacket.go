package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_mysql_FieldPacket struct {
  BuiltinPrototype
}

func NewNodeJS_mysql_FieldPacketPrototype() values.Prototype {
  return &NodeJS_mysql_FieldPacket{newBuiltinPrototype("FieldPacket")}
}

func NewNodeJS_mysql_FieldPacket(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_mysql_FieldPacketPrototype(), ctx)
}

func (p *NodeJS_mysql_FieldPacket) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_mysql_FieldPacket); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_mysql_FieldPacket) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "catalog", "db", "table", "name":
    return s, nil
  case "length", "flags", "type":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_mysql_FieldPacket) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_mysql_FieldPacketPrototype(), ctx), nil
}
