package values

import (
  "strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Array struct {
  content Value
  length int
  ValueData
}

func NewArray(content Value, length int, ctx context.Context) Value {
  if length < 2 {
    panic("pointless array")
  }

  return &Array{content, length, newValueData(ctx)}
}

func (v *Array) TypeName() string {
  return v.content.TypeName() + "[" + strconv.Itoa(v.length) + "]"
}

func (v *Array) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if other, ok := other_.(*Array); ok {
    if other.length == v.length {
      if err := v.content.Check(other.content, ctx); err == nil {
        return nil
      }
    }
  }

  return ctx.NewError("Error: expected " + v.TypeName() + ", got " + other_.TypeName())
}

func (v *Array) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *Array) GetMember(key string, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get member of array")
}

func (v *Array) SetMember(key string, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set member of array")
}

func (v *Array) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  i, _ := idx.LiteralIntValue()

  if i < 0 || i >= v.length {
    return nil, ctx.NewError("Error: index " + strconv.Itoa(i) + " out of range")
  }

  return NewContextValue(v.content, ctx), nil
}

func (v *Array) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  i, _ := idx.LiteralIntValue()

  if i < 0 || i >= v.length {
    return ctx.NewError("Error: index " + strconv.Itoa(i) + " out of range")
  }

  return v.content.Check(arg, ctx)
}

func (v *Array) Length() int {
  return v.length
}
