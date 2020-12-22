package prototypes

import (
  "strings"

  "../values"

  "../../context"
)

type GLEnum struct {
  name string 
  BuiltinPrototype
}

func NewGLEnumPrototype(name string) values.Prototype {
  return &GLEnum{name, newBuiltinPrototype("GLEnum")}
}

func NewGLEnum(ctx context.Context) values.Value {
  return values.NewInstance(NewGLEnumPrototype(""), ctx)
}

func NewNamedGLEnum(name string, ctx context.Context) values.Value {
  return values.NewInstance(NewGLEnumPrototype(name), ctx)
}

func (p *GLEnum) Name() string {
  var b strings.Builder

  b.WriteString("GLEnum")

  if p.name != "" {
    b.WriteString("<")
    b.WriteString(p.name)
    b.WriteString(">")
  }

  return b.String()
}

func (p *GLEnum) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*GLEnum); ok {
    if p.name == "" {
      return nil
    } else if p.name == other.name {
      return nil
    } else {
      return ctx.NewError("Error: expected literal " + p.Name() + ", got " + other.Name()) 
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *GLEnum) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewGLEnumPrototype(""), ctx), nil
}
