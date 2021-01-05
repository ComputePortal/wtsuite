package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SamplerCubeType struct {
  TypeData
}

func NewSamplerCubeType(ctx context.Context) Value {
  return &SamplerCubeType{newTypeData("samplerCube", ctx)}
}

func (v *SamplerCubeType) Check(other Value, ctx context.Context) error {
  instance, _ := v.Instantiate(v.Context())

  return instance.Check(other, ctx)
}

func (v *SamplerCubeType) Instantiate(ctx context.Context) (Value, error) {
  return NewSamplerCube(ctx), nil
}

func (v *SamplerCubeType) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a constructor")
}
