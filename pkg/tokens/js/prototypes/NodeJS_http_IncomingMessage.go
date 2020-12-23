package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type NodeJS_http_IncomingMessage struct {
  BuiltinPrototype
}

func NewNodeJS_http_IncomingMessagePrototype() values.Prototype {
  return &NodeJS_http_IncomingMessage{newBuiltinPrototype("IncomingMessage")}
}

func NewNodeJS_http_IncomingMessage(ctx context.Context) values.Value {
  return values.NewInstance(NewNodeJS_http_IncomingMessagePrototype(), ctx)
}

func (p *NodeJS_http_IncomingMessage) GetParent() (values.Prototype, error) {
  return NewNodeJS_stream_ReadablePrototype(), nil
}

func (p *NodeJS_http_IncomingMessage) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*NodeJS_http_IncomingMessage); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *NodeJS_http_IncomingMessage) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "aborted", "complete":
    return b, nil
  case "headers":
    return NewMapLikeObject(s, ctx), nil
  case "httpVersion", "method", "statusMessage", "url":
    return s, nil
  case "rawHeaders":
    return NewArray(s, ctx), nil
  case "statusCode":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *NodeJS_http_IncomingMessage) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewNodeJS_http_IncomingMessagePrototype(), ctx), nil
}
