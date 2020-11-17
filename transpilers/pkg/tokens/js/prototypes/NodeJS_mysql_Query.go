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

func (p *NodeJS_mysql_Query) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_mysql_QueryPrototype(), ctx), nil
}
