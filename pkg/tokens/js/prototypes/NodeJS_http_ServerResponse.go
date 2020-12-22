package prototypes

import (
  "../values"

  "../../context"
)

type NodeJS_http_ServerResponse struct {
  BuiltinPrototype
}

func NewNodeJS_http_ServerResponsePrototype() values.Prototype {
  return &NodeJS_http_ServerResponse{newBuiltinPrototype("ServerResponse")}
}

func NewNodeJS_http_ServerResponse(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_http_ServerResponsePrototype(), ctx)
}

func (p *NodeJS_http_ServerResponse) GetParent() (values.Prototype, error) {
  return NewNodeJS_EventEmitterPrototype(), nil
}

func (p *NodeJS_http_ServerResponse) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_http_ServerResponse); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_http_ServerResponse) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  o := NewObject(nil, ctx)
  s := NewString(ctx)
  buf := NewNodeJS_Buffer(ctx)
  resp := NewNodeJS_http_ServerResponse(ctx)

  switch key {
  case "statusCode", "statusMessage":
    return nil, ctx.NewError("Error: only a setter")
  case "end":
    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{resp},
      []values.Value{s, resp},
      []values.Value{buf, resp},
      []values.Value{s, s, resp},
      []values.Value{buf, s, resp},
    }, ctx), nil
  case "write":
    callback := values.NewFunction([]values.Value{nil}, ctx)

    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{s, b},
      []values.Value{s, s, b},
      []values.Value{s, callback, b},
      []values.Value{s, s, callback, b},
      []values.Value{buf, b},
      []values.Value{buf, s, b},
      []values.Value{buf, callback, b},
      []values.Value{buf, s, callback, b},
    }, ctx), nil
  case "writeContinue":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "writeHead":
    ss := NewArray(s, ctx)
    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{i, s, o, resp},
      []values.Value{i, o, resp},
      []values.Value{i, s, ss, resp},
      []values.Value{i, ss, resp},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_http_ServerResponse) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "statusCode":
    return i.Check(arg, ctx)
  case "statusMessage":
    return s.Check(arg, ctx)
  default:
    return ctx.NewError("Error: http.ServerResponse." + key + " not setable")
  }
}

func (p *NodeJS_http_ServerResponse) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_http_ServerResponsePrototype(), ctx), nil
}
