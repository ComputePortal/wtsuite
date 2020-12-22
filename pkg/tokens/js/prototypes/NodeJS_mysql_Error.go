package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_mysql_Error struct {
  BuiltinPrototype
}

func NewNodeJS_mysql_ErrorPrototype() values.Prototype {
  return &NodeJS_mysql_Error{newBuiltinPrototype("Error")}
}

func NewNodeJS_mysql_Error(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_mysql_ErrorPrototype(), ctx)
}

func (p *NodeJS_mysql_Error) GetParent() (values.Prototype, error) {
  return NewErrorPrototype(), nil
}

func (p *NodeJS_mysql_Error) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_mysql_Error); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_mysql_Error) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "code", "sqlMessage", "sqlState", "sql":
    return s, nil
  case "errno", "index":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_mysql_Error) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_mysql_ErrorPrototype(), ctx), nil
}
