package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// give a different context to a value
type ContextValue struct {
	val Value
	ctx context.Context
}

func NewContextValue(val Value, ctx context.Context) Value {
  return &ContextValue{UnpackContextValue(val), ctx}
}

func UnpackContextValue(val Value) Value {
	if ctxVal, ok := val.(*ContextValue); ok {
		return UnpackContextValue(ctxVal.val)
	} else {
		return val
	}
}

func (v *ContextValue) Context() context.Context {
  return v.ctx
}

func (v *ContextValue) TypeName() string {
  return v.val.TypeName()
}

func (v *ContextValue) Check(v_ Value, ctx context.Context) error {
  return v.val.Check(v_, ctx)
}

func (v *ContextValue) Instantiate(ctx context.Context) (Value, error) {
	return v.val.Instantiate(ctx)
}

func (v *ContextValue) EvalFunction(args []Value, ctx context.Context) (Value, error) {
	return v.val.EvalFunction(args, ctx)
}

func (v *ContextValue) GetMember(key string, ctx context.Context) (Value, error) {
	return v.val.GetMember(key, ctx)
}

func (v *ContextValue) SetMember(key string, arg Value, ctx context.Context) error {
	return v.val.SetMember(key, arg, ctx)
}

func (v *ContextValue) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
	return v.val.GetIndex(idx, ctx)
}

func (v *ContextValue) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
	return v.val.SetIndex(idx, arg, ctx)
}

func (v *ContextValue) LiteralIntValue() (int, bool) {
  return v.val.LiteralIntValue()
}

func (v *ContextValue) Length() int {
  return v.val.Length()
}
