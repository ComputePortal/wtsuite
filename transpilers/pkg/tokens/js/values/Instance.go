package values

import (
	"../../context"
)

type Instance struct {
	interf Interface

	ValueData
}

func newInstance(interf Interface, ctx context.Context) Instance {
  return Instance{interf, ValueData{ctx}}
}
func NewInstance(interf Interface, ctx context.Context) Value {
  inst := newInstance(interf, ctx)
  return &inst
}

func (v *Instance) TypeName() string {
	return v.interf.Name()
}

func (v *Instance) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAll(other_) {
    return nil
  } else if other, ok := other_.(*Instance); ok {
    // first match the interface
    if err := v.interf.Check(other.interf, ctx); err != nil {
      return err
    }

    return nil
  } else {
    return ctx.NewError("Error: not an instance")
  }
}

func (v *Instance) GetInterface() Interface {
  return v.interf
}

func (v *Instance) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't construct an instance")
}

func (v *Instance) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't call an instance")
}

func (v *Instance) GetMember(key string, includePrivate bool,
	ctx context.Context) (Value, error) {
  return v.interf.GetInstanceMember(key, includePrivate, ctx)
}

func (v *Instance) SetMember(key string, includePrivate bool, arg Value,
	ctx context.Context) error {
  return v.interf.SetInstanceMember(key, includePrivate, arg, ctx)
}

func IsInstance(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  switch v_.(type) {
  case *LiteralStringInstance:
    return true
  case *LiteralBooleanInstance:
    return true
  case *LiteralIntInstance:
    return true
  case *Instance:
    return true
  default:
    return false
  }
}

func IsLiteral(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  switch v_.(type) {
  case *LiteralStringInstance:
    return true
  case *LiteralBooleanInstance:
    return true
  case *LiteralIntInstance:
    return true
  default:
    return false
  }
}
