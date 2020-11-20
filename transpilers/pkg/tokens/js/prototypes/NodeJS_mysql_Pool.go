package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_mysql_Pool struct {
  BuiltinPrototype
}

func NewNodeJS_mysql_PoolPrototype() values.Prototype {
  return &NodeJS_mysql_Pool{newBuiltinPrototype("Pool")}
}

func NewNodeJS_mysql_Pool(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_mysql_PoolPrototype(), ctx)
}

func (p *NodeJS_mysql_Pool) GetParent() (values.Prototype, error) {
  return NewNodeJS_mysql_ConnectionPrototype(), nil
}

func (p *NodeJS_mysql_Pool) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_mysql_Pool); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_mysql_Pool) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "getConnection":
    callback := values.NewFunction([]values.Value{
      NewNodeJS_mysql_Error(ctx),
      NewNodeJS_mysql_Connection(ctx),
      nil,
    }, ctx)

    return values.NewFunction([]values.Value{callback, nil}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_mysql_Pool) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_mysql_PoolPrototype(), ctx), nil
}
