package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Struct struct {
  t Structable
  ValueData
}

func NewStruct(t Structable, ctx context.Context) Value {
  return &Struct{t, newValueData(ctx)}
}

func (v *Struct) GetStructable() Structable {
  return v.t
}

func (v *Struct) TypeName() string {
  return "struct " + v.t.Name()
}

func (v *Struct) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if other, ok := other_.(*Struct); ok {
    if other.t == v.t {
      return nil
    }
  }

  return ctx.NewError("Error: expected " + v.TypeName() + ", got " + other_.TypeName())
}

func (v *Struct) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *Struct) GetMember(key string, ctx context.Context) (Value, error) {
  return v.t.GetMember(key, ctx) 
}

func (v *Struct) SetMember(key string, arg Value, ctx context.Context) error {
  return v.t.SetMember(key, arg, ctx) 
}

func (v *Struct) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get index of struct")
}

func (v *Struct) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set index of struct")
}

func (v *Struct) Length() int {
  return 1
}

func IsStruct(v Value) bool {
  v = UnpackContextValue(v)

  _, ok := v.(*Struct)
  return ok
}

func AssertStruct(v Value) (*Struct, error) {
  errCtx := v.Context()
  v = UnpackContextValue(v)

  if vStruct, ok := v.(*Struct); ok {
    return vStruct, nil
  } else {
    return nil, errCtx.NewError("Error: not a struct")
  }
}
