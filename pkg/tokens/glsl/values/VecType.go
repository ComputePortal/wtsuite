package values

import (
  "strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type VecType struct {
  compType string // float, int or bool
  n int // 2, 3 or 4
  TypeData
}

func vecTypeName(compType string, n int) string {
  typeName := "vec" + strconv.Itoa(n)

  switch compType {
  case "bool":
    typeName = "b" + typeName
  case "int":
    typeName = "i" + typeName
  case "float":
    // ok
  default:
    panic("unhandled")
  }

  if n < 2 || n > 4 {
    panic("unhandled")
  }

  return typeName
}

func NewVecType(compType string, n int, ctx context.Context) Value {
  typeName := vecTypeName(compType, n)

  return &VecType{compType, n, newTypeData(typeName, ctx)}
}

func (v *VecType) Check(other Value, ctx context.Context) error {
  instance, _ := v.Instantiate(v.Context())

  return instance.Check(other, ctx)
}

func (v *VecType) Instantiate(ctx context.Context) (Value, error) {
  return NewVec(v.compType, v.n, ctx), nil
}

// can't vecs also take other vecs as arguments?
func (v *VecType) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  if len(args) != v.n {
    return nil, ctx.NewError("Error: expected " + strconv.Itoa(v.n) + " arguments")
  }

  for i := 0; i < v.n; i++ {
    argCtx := args[i].Context()
    arg := UnpackContextValue(args[i])

    if _, ok := arg.(*Scalar); !ok {
      errCtx := argCtx
      return nil, errCtx.NewError("Error: expected scalar argument")
    }
  }

  return v.Instantiate(ctx)
}
