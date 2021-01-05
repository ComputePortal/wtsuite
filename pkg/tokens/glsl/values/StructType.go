package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type StructType struct {
  t Structable
  TypeData
}

func NewStructType(t Structable, ctx context.Context) Value {
  return &StructType{t, newTypeData("struct " + t.Name(), ctx)}
}

func (v *StructType) Check(other Value, ctx context.Context) error {
  instance, _ := v.Instantiate(v.Context())

  return instance.Check(other, ctx)
}

func (v *StructType) Instantiate(ctx context.Context) (Value, error) {
  return NewStruct(v.t, ctx), nil
}

func (v *StructType) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  if err := v.t.CheckConstruction(args, ctx); err != nil {
    return nil, err
  }

  return v.Instantiate(ctx)
}
