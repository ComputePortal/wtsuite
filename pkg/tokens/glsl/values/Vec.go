package values

import (
  "strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Vec struct {
  compName string
  n int
  ValueData
}

func NewVec(compName string, n int, ctx context.Context) Value {
  return &Vec{compName, n, newValueData(ctx)}
}

func NewVec2(ctx context.Context) Value {
  return NewVec("float", 2, ctx)
}

func NewVec3(ctx context.Context) Value {
  return NewVec("float", 3, ctx)
}

func NewVec4(ctx context.Context) Value {
  return NewVec("float", 4, ctx)
}

func NewIVec2(ctx context.Context) Value {
  return NewVec("int", 2, ctx)
}

func NewIVec3(ctx context.Context) Value {
  return NewVec("int", 3, ctx)
}

func NewIVec4(ctx context.Context) Value {
  return NewVec("int", 4, ctx)
}

func NewBVec2(ctx context.Context) Value {
  return NewVec("bool", 2, ctx)
}

func NewBVec3(ctx context.Context) Value {
  return NewVec("bool", 3, ctx)
}

func NewBVec4(ctx context.Context) Value {
  return NewVec("bool", 4, ctx)
}

func (v *Vec) TypeName() string {
  return vecTypeName(v.compName, v.n)
}

func (v *Vec) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if other, ok := other_.(*Vec); ok {
    if other.compName == v.compName && other.n == v.n {
      return nil
    }
  }

  return ctx.NewError("Error: expected " + v.TypeName() + ", got " + other_.TypeName())
}

func (v *Vec) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *Vec) assertValidSwizzleKey(key string, ctx context.Context) error {
  if len(key) > v.n {
    return ctx.NewError("Error: " + v.TypeName() + "." + key + " not found")
  }

  for i := 0; i < len(key); i++ {
    c := key[i]

    if c == 'x' || c == 'y' {
      // always ok
      continue
    } else if c == 'z' && v.n > 2 {
      continue
    } else if c == 'w' && v.n > 3 {
      continue
    } else {
      return ctx.NewError("Error: " + v.TypeName() + "." + key + " not found")
    }
  }

  return nil
}

func (v *Vec) GetMember(key string, ctx context.Context) (Value, error) {
  if err := v.assertValidSwizzleKey(key, ctx); err != nil {
    return nil, err
  }

  switch len(key) {
  case 1:
    return NewScalar(v.compName, ctx), nil
  case 2:
    return NewVec(v.compName, 2, ctx), nil
  case 3:
    return NewVec(v.compName, 3, ctx), nil
  case 4:
    return NewVec(v.compName, 4, ctx), nil
  default:
    return nil, ctx.NewError("Error: " + v.TypeName() + "." + key + " not found")
  }
}

func (v *Vec) SetMember(key string, arg Value, ctx context.Context) error {
  if err := v.assertValidSwizzleKey(key, ctx); err != nil {
    return err
  }

  switch len(key) {
  case 1:
    check := NewScalar(v.compName, ctx)
    return check.Check(arg, ctx)
  case 2:
    check := NewVec(v.compName, 2, ctx)
    return check.Check(arg, ctx)
  case 3:
    check := NewVec(v.compName, 3, ctx)
    return check.Check(arg, ctx)
  case 4:
    check := NewVec(v.compName, 4, ctx)
    return check.Check(arg, ctx)
  default:
    return ctx.NewError("Error: " + v.TypeName() + "." + key + " not found")
  }
}

func (v *Vec) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  i, _ := idx.LiteralIntValue()

  if i < 0 || i >= v.n {
    return nil, ctx.NewError("Error: index " + strconv.Itoa(i) + " out of range")
  }

  return NewScalar(v.compName, ctx), nil
}

func (v *Vec) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  i, _ := idx.LiteralIntValue()

  if i < 0 || i >= v.n {
    return ctx.NewError("Error: index " + strconv.Itoa(i) + " out of range")
  }

  check := NewScalar(v.compName, ctx)

  return check.Check(arg, ctx)
}

func (v *Vec) Length() int {
  return v.n
}

func IsVec(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  _, ok := v_.(*Vec)
  return ok
}
