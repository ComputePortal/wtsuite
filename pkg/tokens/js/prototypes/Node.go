package prototypes

import (
  "../values"

  "../../context"
)

type Node struct {
  BuiltinPrototype
}

func NewNodePrototype() values.Prototype {
  return &Node{newBuiltinPrototype("Node")}
}

func NewNode(ctx context.Context) values.Value {
  return values.NewInstance(NewNodePrototype(), ctx)
}

func (p *Node) GetParent() (values.Prototype, error) {
  return NewEventTargetPrototype(), nil
}

func (p *Node) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Node); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Node) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  n := NewNode(ctx)

  switch key {
  case "appendChild", "removeChild":
    return values.NewMethodLikeFunction([]values.Value{n, n}, ctx), nil
  case "contains":
    return values.NewFunction([]values.Value{n, b}, ctx), nil
  case "firstChild", "lastChild", "parentNode":
    return n, nil
  case "insertBefore", "replaceChild":
    return values.NewMethodLikeFunction([]values.Value{n, n, n}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Node) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodePrototype(), ctx), nil
}
