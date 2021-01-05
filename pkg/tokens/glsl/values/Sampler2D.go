package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Sampler2D struct {
  ValueData
}

func NewSampler2D(ctx context.Context) Value {
  return &Sampler2D{newValueData(ctx)}
}

func (v *Sampler2D) TypeName() string {
  return "sampler2D"
}

func IsSampler2D(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  _, ok := v_.(*Sampler2D)
  return ok
}

func (v *Sampler2D) Check(other_ Value, ctx context.Context) error {
  if IsSampler2D(other_) {
    return nil
  } else {
    return ctx.NewError("Error: expected " + v.TypeName() + ", got " + other_.TypeName())
  }
}

func (v *Sampler2D) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *Sampler2D) GetMember(key string, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get member of sampler2D")
}

func (v *Sampler2D) SetMember(key string, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set member of sampler2D")
}

func (v *Sampler2D) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get index of sampler2D")
}

func (v *Sampler2D) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set index of sampler2D")
}

func (v *Sampler2D) LiteralIntValue() (int, bool) {
  return 0, false
}

func (v *Sampler2D) Length() int {
  return 1
}
