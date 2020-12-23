package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type RPCServer struct {
  BuiltinPrototype
}

func NewRPCServerPrototype() values.Prototype {
  return &RPCServer{newBuiltinPrototype("RPCServer")}
}

func NewRPCServer(ctx context.Context) values.Value {
  return values.NewInstance(NewRPCServerPrototype(), ctx)
}

func (p *RPCServer) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*RPCServer); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *RPCServer) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  s := NewString(ctx)

  switch key {
  case "handle":
    return values.NewFunction([]values.Value{s, NewPromise(s, ctx)}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *RPCServer) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewRPCServerPrototype(), ctx), nil
}
