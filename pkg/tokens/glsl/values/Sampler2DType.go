package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Sampler2DType struct {
  TypeData
}

func NewSampler2DType(ctx context.Context) Value {
  return &Sampler2DType{newTypeData("sampler2D", ctx)}
}

func (v *Sampler2DType) Check(other Value, ctx context.Context) error {
  instance, _ := v.Instantiate(v.Context())

  return instance.Check(other, ctx)
}

func (v *Sampler2DType) Instantiate(ctx context.Context) (Value, error) {
  return NewSampler2D(ctx), nil
}

func (v *Sampler2DType) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a constructor")
}
