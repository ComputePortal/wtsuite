package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Function struct {
  fn Callable
  ValueData
}

func NewFunction(fn Callable, ctx context.Context) Value {
  return &Function{fn, newValueData(ctx)}
}

func (v *Function) TypeName() string {
  return "function"
}

// same interface isn't good enough, needs to be a
func (v *Function) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if other, ok := other_.(*Function); ok {
    if other.fn == v.fn {
      return nil
    }

    return ctx.NewError("Error: function has different interface")
  }

  return ctx.NewError("Error: expected function, got " + other_.TypeName())
}

func (v *Function) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return v.fn.EvalCall(args, ctx)
}

func (v *Function) GetMember(key string, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get member of function")
}

func (v *Function) SetMember(key string, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set member of function")
}

func (v *Function) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get index of function")
}

func (v *Function) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set index of function")
}

func (v *Function) LiteralIntValue() (int, bool) {
  return 0, false
}

func (v *Function) Length() int {
  return 1
}

func AssertMainFunction(v Value) error {
  ctx := v.Context()

  if retVal, err := v.EvalFunction([]Value{}, v.Context()); err != nil {
    return ctx.NewError("Error: invalid void main() function")
  } else if retVal != nil {
    return ctx.NewError("Error: return value not void")
  }

  return nil
}
