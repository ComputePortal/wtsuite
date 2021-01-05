package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ScalarType struct {
  TypeData
}

func NewScalarType(name string, ctx context.Context) Value {
  return &ScalarType{newTypeData(name, ctx)}
}

func (v *ScalarType) Check(other Value, ctx context.Context) error {
  instance, _ := v.Instantiate(v.Context())

  return instance.Check(other, ctx)
}

func (v *ScalarType) Instantiate(ctx context.Context) (Value, error) {
  return NewScalar(v.name, ctx), nil
}

func (v *ScalarType) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  if len(args) != 1 {
    return nil, ctx.NewError("Error: expected 1 argument")
  }

  arg0 := UnpackContextValue(args[0])

  if _, ok := arg0.(*Scalar); !ok {
    errCtx := arg0.Context()
    return nil, errCtx.NewError("Error: expected scalar argument")
  }

  return v.Instantiate(ctx)
}
