package prototypes

import (
  "../values"

  "../../context"
)

type WebGLProgram struct {
  BuiltinPrototype
}

func NewWebGLProgramPrototype() values.Prototype {
  return &WebGLProgram{newBuiltinPrototype("WebGLProgram")}
}

func NewWebGLProgram(ctx context.Context) values.Value {
  return values.NewInstance(NewWebGLProgramPrototype(), ctx)
}

func (p *WebGLProgram) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebGLProgram); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *WebGLProgram) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLProgramPrototype(), ctx), nil
}
