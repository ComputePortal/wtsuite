package prototypes

import (
  "../values"

  "../../context"
)

type TextMetrics struct {
  BuiltinPrototype
}

func NewTextMetricsPrototype() values.Prototype {
  return &TextMetrics{newBuiltinPrototype("TextMetrics")}
}

func NewTextMetrics(ctx context.Context) values.Value {
  return values.NewInstance(NewTextMetricsPrototype(), ctx)
}

func (p *TextMetrics) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  f := NewNumber(ctx)

  switch key {
  case "width":
    return f, nil
  default:
    return nil, nil
  }
}

func (p *TextMetrics) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewUnconstructableClass(NewTextMetricsPrototype(), ctx), nil
}
