package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Text struct {
  BuiltinPrototype
}

func NewTextPrototype() values.Prototype {
  return &Text{newBuiltinPrototype("Text")}
}

func NewText(ctx context.Context) values.Value {
  return values.NewInstance(NewTextPrototype(), ctx)
}

func (p *Text) GetParent() (values.Prototype, error) {
  return NewNodePrototype(), nil
}

func (p *Text) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Text); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Text) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewTextPrototype(), ctx), nil
}
