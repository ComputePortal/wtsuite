package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SamplerCube struct {
  ValueData
}

func NewSamplerCube(ctx context.Context) Value {
  return &SamplerCube{newValueData(ctx)}
}

func (v *SamplerCube) TypeName() string {
  return "samplerCube"
}

func IsSamplerCube(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  _, ok := v_.(*SamplerCube)
  return ok
}

func (v *SamplerCube) Check(other_ Value, ctx context.Context) error {
  if IsSamplerCube(other_) {
    return nil
  } else {
    return ctx.NewError("Error: expected " + v.TypeName() + ", got " + other_.TypeName())
  }
}

func (v *SamplerCube) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *SamplerCube) GetMember(key string, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get member of samplerCube")
}

func (v *SamplerCube) SetMember(key string, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set member of samplerCube")
}

func (v *SamplerCube) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get index of samplerCube")
}

func (v *SamplerCube) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set index of samplerCube")
}

func (v *SamplerCube) LiteralIntValue() (int, bool) {
  return 0, false
}

func (v *SamplerCube) Length() int {
  return 1
}
