package values

import (
  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type FunctionOverloads struct {
  args [][]Value // each list is an overload
}

func NewBuiltinFunction(args [][]Value, ctx context.Context) Value {
  return &Function{&FunctionOverloads{args}, newValueData(ctx)}
}

func (fo *FunctionOverloads) EvalCall(args []Value, ctx context.Context) (Value, error) {
  i, err := checkAnyOverload(fo.args, args, ctx)
  if err != nil {
    return nil, err
  }

  overload := fo.args[i]

  return overload[len(overload)-1], nil
}

func checkAnyOverload(overloads [][]Value, ts []Value, ctx context.Context) (int, error) {
  for i, overload_ := range overloads {
    overload := overload_[0:len(overload_)-1] // cut off the return value

    if len(overload) == len(ts) {
      ok := true
      for j, arg := range overload {
        if err := arg.Check(ts[j], ts[j].Context()); err != nil {
          if len(overloads) == 1 {
            return 0, err
          }

          ok = false
          break
        }
      }

      if ok {
        return i, nil
      }
    }
  }

  return 0, ctx.NewError("Error: expected other arg types")
}

func NewOneToOneFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f},
    []Value{v2, v2},
    []Value{v3, v3},
    []Value{v4, v4},
  }, ctx)
}

func NewOneOrTwoToOneFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f},
    []Value{f, f, f},
    []Value{v2, v2},
    []Value{v2, v2, v2},
    []Value{v3, v3},
    []Value{v3, v3, v3},
    []Value{v4, v4},
    []Value{v4, v4, v4},
  }, ctx)
}

func NewTwoToOneFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f},
    []Value{v2, v2, v2},
    []Value{v3, v3, v3},
    []Value{v4, v4, v4},
  }, ctx)
}

func NewThreeToOneFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f, f},
    []Value{v2, v2, v2, v2},
    []Value{v3, v3, v3, v3},
    []Value{v4, v4, v4, v4},
  }, ctx)
}

func NewMinMaxFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f},
    []Value{v2, v2, v2},
    []Value{v2, f, v2},
    []Value{v3, v3, v3},
    []Value{v3, f, v3},
    []Value{v4, v4, v4},
    []Value{v4, f, v4},
  }, ctx)
}

func NewClampFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f, f},
    []Value{v2, v2, v2, v2},
    []Value{v3, v3, v3, v3},
    []Value{v4, v4, v4, v4},
    []Value{v2, f, f, v2},
    []Value{v3, f, f, v3},
    []Value{v4, f, f, v4},
  }, ctx)
}

func NewMixFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f, f},
    []Value{v2, v2, v2, v2},
    []Value{v3, v3, v3, v3},
    []Value{v4, v4, v4, v4},
    []Value{v2, v2, f, v2},
    []Value{v3, v3, f, v3},
    []Value{v4, v4, f, v4},
  }, ctx)
}

func NewStepFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f},
    []Value{v2, v2, v2},
    []Value{v3, v3, v3},
    []Value{v4, v4, v4},
    []Value{f, v2, v2},
    []Value{f, v3, v3},
    []Value{f, v4, v4},
  }, ctx)
}

func NewSmoothStepFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f, f},
    []Value{v2, v2, v2, v2},
    []Value{v3, v3, v3, v3},
    []Value{v4, v4, v4, v4},
    []Value{f, f, v2, v2},
    []Value{f, f, v3, v3},
    []Value{f, f, v4, v4},
  }, ctx)
}

func NewLengthFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f},
    []Value{v2, f},
    []Value{v3, f},
    []Value{v4, f},
  }, ctx)
}

func NewDotFunction(ctx context.Context) Value {
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{f, f, f},
    []Value{v2, v2, f},
    []Value{v3, v3, f},
    []Value{v4, v4, f},
  }, ctx)
}

func NewCrossFunction(ctx context.Context) Value {
  v3 := NewVec3(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{v3, v3, v3},
  }, ctx)
}

func NewCompareFunction(ctx context.Context) Value {
  v2 := NewVec2(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  i2 := NewIVec2(ctx)
  i3 := NewIVec3(ctx)
  i4 := NewIVec4(ctx)

  b2 := NewBVec2(ctx)
  b3 := NewBVec3(ctx)
  b4 := NewBVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{v2, v2, b2},
    []Value{v3, v3, b3},
    []Value{v4, v4, b4},
    []Value{i2, i2, b2},
    []Value{i3, i3, b3},
    []Value{i4, i4, b4},
  }, ctx)
}

func NewAnyAllFunction(ctx context.Context) Value {
  b := NewBool(ctx)
  b2 := NewBVec2(ctx)
  b3 := NewBVec3(ctx)
  b4 := NewBVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{b2, b},
    []Value{b3, b},
    []Value{b4, b},
  }, ctx)
}

func NewNotFunction(ctx context.Context) Value {
  b2 := NewBVec2(ctx)
  b3 := NewBVec3(ctx)
  b4 := NewBVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{b2, b2},
    []Value{b3, b3},
    []Value{b4, b4},
  }, ctx)
}

func NewTexture2DFunction(ctx context.Context) Value {
  s2D := NewSampler2D(ctx)
  f := NewFloat(ctx)
  v2 := NewVec2(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{s2D, v2, v4},
    []Value{s2D, v2, f, v4},
  }, ctx)
}

func NewTextureCubeFunction(ctx context.Context) Value {
  sC := NewSamplerCube(ctx)
  f := NewFloat(ctx)
  v3 := NewVec3(ctx)
  v4 := NewVec4(ctx)

  return NewBuiltinFunction([][]Value{
    []Value{sC, v3, v4},
    []Value{sC, v3, f, v4},
  }, ctx)
}
