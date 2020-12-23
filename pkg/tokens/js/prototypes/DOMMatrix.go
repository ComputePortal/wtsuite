package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type DOMMatrix struct {
  BuiltinPrototype
}

func NewDOMMatrixPrototype() values.Prototype {
  return &DOMMatrix{newBuiltinPrototype("DOMMatrix")}
}

func NewDOMMatrix(ctx context.Context) values.Value {
  return values.NewInstance(NewDOMMatrixPrototype(), ctx)
}

func (p *DOMMatrix) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*DOMMatrix); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *DOMMatrix) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewDOMMatrixPrototype(), ctx), nil
}
