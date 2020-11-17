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

func (p *WebGLProgram) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLProgramPrototype(), ctx), nil
}
