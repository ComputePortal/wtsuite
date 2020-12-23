package prototypes

import (
  "strings"

  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type IDBRequest struct {
  content values.Value // if nil, then any

  BuiltinPrototype
}

func NewIDBRequestPrototype(content values.Value) values.Prototype {
  return &IDBRequest{content, newBuiltinPrototype("IDBRequest")}
}

func NewIDBRequest(content values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewIDBRequestPrototype(content), ctx)
}

func NewEmptyIDBRequest(ctx context.Context) values.Value {
  return NewIDBRequest(nil, ctx)
}

func (p *IDBRequest) Name() string {
  var b strings.Builder

  b.WriteString("IDBRequest")

  if p.content != nil {
    b.WriteString("<")
    b.WriteString(p.content.TypeName())
    b.WriteString(">")
  }

  return b.String()
}

func (p *IDBRequest) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*IDBRequest); ok {
    if p.content == nil {
      return nil
    } else if other.content == nil {
      return ctx.NewError("Error: expected IDBRequest<" + p.content.TypeName() + ">, got IDBRequest<any>")
    } else if p.content.Check(other.content, ctx) != nil {
      return ctx.NewError("Error: expected IDBRequest<" + p.content.TypeName() + ">, got IDBRequest<" + other.content.TypeName() + ">")
    } else {
      return nil
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *IDBRequest) getContentValue() values.Value {
  if p.content == nil {
    return values.NewAny(context.NewDummyContext())
  } else {
    return p.content
  }
}

func (p *IDBRequest) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  content := values.NewContextValue(p.getContentValue(), ctx)

  switch key {
  case "onerror", "onsuccess":
    return nil, ctx.NewError("Error: is a setter only")
  case "result":
    return content, nil
  default:
    return nil, nil
  }
}

func (p *IDBRequest) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  callback := values.NewFunction([]values.Value{NewEvent(NewIDBRequest(values.NewAny(ctx), ctx), ctx), nil}, ctx)

  switch key {
  case "onerror", "onsuccess":
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: IDBRequest." + key + " not setable")
  }
}

func (p *IDBRequest) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewIDBRequestPrototype(values.NewAny(ctx)), ctx), nil
}
