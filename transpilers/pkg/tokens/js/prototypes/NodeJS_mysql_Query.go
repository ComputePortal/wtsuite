package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_mysql_Query struct {
  BuiltinPrototype
}

func NewNodeJS_mysql_QueryPrototype() values.Prototype {
  return &NodeJS_mysql_Query{newBuiltinPrototype("Query")}
}

func NewNodeJS_mysql_Query(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_mysql_QueryPrototype(), ctx)
}

func (p *NodeJS_mysql_Query) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_mysql_Query) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_mysql_Query); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_mysql_Query) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_mysql_QueryPrototype(), ctx), nil
}
