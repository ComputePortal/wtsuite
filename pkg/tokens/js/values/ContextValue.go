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
	if ctxVal, ok := val.(*ContextValue); ok {
		return NewContextValue(ctxVal.val, ctx)
	} else {
		return &ContextValue{val, ctx}
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

func (v *ContextValue) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
	return v.val.EvalConstructor(args, ctx)
}

func (v *ContextValue) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
	return v.val.EvalFunction(args, preferMethod, ctx)
}

func (v *ContextValue) GetMember(key string, includePrivate bool, 
  ctx context.Context) (Value, error) {
	return v.val.GetMember(key, includePrivate, ctx)
}

func (v *ContextValue) SetMember(key string, includePrivate bool, arg Value,
	ctx context.Context) error {
	return v.val.SetMember(key, includePrivate, arg, ctx)
}

func (v *ContextValue) LiteralBooleanValue() (bool, bool) {
  return v.val.LiteralBooleanValue()
}

func (v *ContextValue) LiteralIntValue() (int, bool) {
  return v.val.LiteralIntValue()
}

func (v *ContextValue) LiteralStringValue() (string, bool) {
  return v.val.LiteralStringValue()
}

func UnpackContextValue(val Value) Value {
  var res Value = val

  for true {
    br := false
    switch res_ := res.(type) {
    case *ContextValue:
      res = res_.val
    case *This:
      res = res_.this
    default:
      br = true
    }

    if br {
      break
    }
  }

  return res
}
