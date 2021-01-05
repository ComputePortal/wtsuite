package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
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
  if other_ == nil {
    panic("other_ can't be nil")
  }

  other_ = UnpackContextValue(other_)

  switch other := other_.(type) {
  case *Instance: 
    // first match the interface
    if err := v.interf.Check(other.interf, ctx); err != nil {
      return err
    }

    return nil
  case *LiteralIntInstance: 
    // first match the interface
    if err := v.interf.Check(other.interf, ctx); err != nil {
      return err
    }

    return nil
  case *LiteralBooleanInstance: 
    // first match the interface
    if err := v.interf.Check(other.interf, ctx); err != nil {
      return err
    }

    return nil
  case *LiteralStringInstance:
    // first match the interface
    if err := v.interf.Check(other.interf, ctx); err != nil {
      return err
    }

    return nil
  default:
    if IsAny(other_) {
      return nil
    } else {
      err := ctx.NewError("Error: have " + other_.TypeName() + ", want " + v.TypeName())

      return err
    }
  }
}

func (v *Instance) GetInterface() Interface {
  return v.interf
}

func (v *Instance) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a constructor")
}

func (v *Instance) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't call an instance")
}

func FindInstanceMemberInterface(interf_ Interface, key string, includePrivate bool, ctx context.Context) (Interface, error) {
  interf := interf_

  for true {
    res, err := interf.GetInstanceMember(key, includePrivate, ctx)
    if err != nil || res != nil {
      return interf, nil
    } 

    if proto, ok := interf.(Prototype); ok {
      parentProto, err := proto.GetParent()
      if err != nil {
        return nil, err
      }

      if parentProto != nil {
        interf = parentProto
      } else {
        break
      }
    } else {
      break
    }
  }

  return nil, ctx.NewError("Error: " + interf_.Name() + "." + key + " not found")
}

func (v *Instance) GetMember(key string, includePrivate bool,
	ctx context.Context) (Value, error) {
  interf, err := FindInstanceMemberInterface(v.interf, key, includePrivate, ctx)
  if err != nil {
    return nil, err
  }

  return interf.GetInstanceMember(key, includePrivate, ctx)
}

func (v *Instance) SetMember(key string, includePrivate bool, arg Value,
	ctx context.Context) error {
  interf, err := FindInstanceMemberInterface(v.interf, key, includePrivate, ctx)
  if err != nil {
    return err
  }

  return interf.SetInstanceMember(key, includePrivate, arg, ctx)
}

func AssertInstance(v_ Value) (*Instance, error) {
  errCtx := v_.Context()
  v_ = UnpackContextValue(v_)

  if v, ok := v_.(*Instance); ok {
    return v, nil
  } 

  return nil, errCtx.NewError("Error: expected an instance, got " + v_.TypeName())
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
  case *Any:
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

func RemoveLiteralness(v Value) Value {
  // literalness must be removed though
  interf := GetInterface(v)

  if interf != nil {
    // potentially literal
    v = NewInstance(interf, v.Context())
  }

  return v
}
