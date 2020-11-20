package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_http_Server struct {
  BuiltinPrototype
}

func NewNodeJS_http_ServerPrototype() values.Prototype {
  return &NodeJS_http_Server{newBuiltinPrototype("Server")}
}

func NewNodeJS_http_Server(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_http_ServerPrototype(), ctx)
}

func (p *NodeJS_http_Server) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_http_Server) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_http_Server); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_http_Server) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "addListener":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewLiteralString("request", ctx), values.NewFunction([]values.Value{
        NewNodeJS_http_IncomingMessage(ctx), NewNodeJS_http_ServerResponse(ctx), nil,
      }, ctx), nil},
      []values.Value{s, values.NewFunction([]values.Value{nil}, ctx), nil},
    }, ctx), nil
  case "listen":
    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{nil},
      []values.Value{s, nil},
      []values.Value{i, nil},
      []values.Value{i, s, nil},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_http_Server) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_http_ServerPrototype(), ctx), nil
}
