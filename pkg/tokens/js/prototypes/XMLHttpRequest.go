package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type XMLHttpRequest struct {
  BuiltinPrototype
}

func NewXMLHttpRequestPrototype() values.Prototype {
  return &XMLHttpRequest{newBuiltinPrototype("XMLHttpRequest")}
}

func NewXMLHttpRequest(ctx context.Context) values.Value {
  return values.NewInstance(NewXMLHttpRequestPrototype(), ctx)
}

func (p *XMLHttpRequest) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*XMLHttpRequest); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *XMLHttpRequest) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "onload", "onerror":
    return nil, ctx.NewError("Error: only a setter")
  case "open", "setRequestHeader":
    return values.NewFunction([]values.Value{s, s, nil}, ctx), nil
  case "status":
    return i, nil
  case "responseText":
    return s, nil
  default:
    return nil, nil
  }
}

func (p *XMLHttpRequest) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  switch key {
  case "onload", "onerror":
    callback := values.NewFunction([]values.Value{nil}, ctx)
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: XMLHttpRequest." + key + " not setable")
  }
}

func (p *XMLHttpRequest) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewClass(
    [][]values.Value{
      []values.Value{},
    }, NewXMLHttpRequestPrototype(), ctx), nil
}
