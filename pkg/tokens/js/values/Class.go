package values

import (
  "fmt"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Class struct {
  args [][]Value // various constructor overloads, can be nil in generic case

  fn func(args []Value, ctx_ context.Context) (Interface, error)

  interf Interface

	ValueData
}

func NewClass(args [][]Value, interf Interface, ctx context.Context) *Class {
	return &Class{args, nil, interf, ValueData{ctx}}
}

func NewCustomClass(args [][]Value, fn func(args []Value, ctx_ context.Context) (Interface, error), ctx context.Context) *Class {
  return &Class{args, fn, nil, ValueData{ctx}}
}

func NewUnconstructableClass(interf Interface, ctx context.Context) *Class {
  return &Class{[][]Value{[]Value{}}, func(args []Value, ctx_ context.Context) (Interface, error) {
    if args == nil {
      return interf, nil
    } else {
      return nil, ctx_.NewError("Error: doesn't have a constructor")
    }
  }, nil, ValueData{ctx}}
}

func (v *Class) GetConstructorArgs() [][]Value {
  return v.args
}

func (v *Class) getInterface() Interface {
  if v.interf == nil {
    if v.fn != nil {
      interf, err := v.fn(nil, context.NewDummyContext())
      if err != nil {
        panic(err)
      }

      return interf
    } else {
      return nil
    }
  } else {
    return v.interf
  }
}

func (v *Class) getPrototype() Prototype {
  interf := v.getInterface()

  if interf == nil {
    return nil
  } else if proto, ok := interf.(Prototype); ok {
    return proto
  } else {
    return nil
  }
}

func (v *Class) TypeName() string {
  var b strings.Builder

  b.WriteString("class")

  // TODO: how to print overloads?
  if v.args != nil {
    b.WriteString("<")
 
    if len(v.args) == 1 {
      for _, arg := range v.args[0] {
        b.WriteString(arg.TypeName())
        b.WriteString(",")
      }
    } else {
      b.WriteString(fmt.Sprintf("%d overloads", len(v.args)))
      b.WriteString(",")
    }

    interf := v.getInterface()
    b.WriteString(interf.Name())

    b.WriteString(">")
  }

  return b.String()
}

// maybe it is a little silly that GetType always needs to be called
func (v *Class) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAny(other_) {
    return nil
  } else if other, ok := other_.(*Class); ok {
    if err := checkAllOverloads(v.args, other.args, ctx); err != nil {
      return err
    }

    interf := v.getInterface()
    otherInterf := other.getInterface()
    if interf == nil {
      return nil
    } else if otherInterf == nil {
      return ctx.NewError("Error: unspecified interface")
    } else {
      return interf.Check(otherInterf, ctx)
    }
  } else {
    return ctx.NewError("Error: not a class")
  }
}

func (v *Class) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return v.evalConstructor(args, ctx, false)
}

func (v *Class) evalConstructor(args []Value, ctx context.Context, allowAbstract bool) (Value, error) {
  if args != nil {
    if _, err := checkAnyOverload(v.args, args, ctx); err != nil {

      for _, overload := range v.args {
        n := len(overload)
        if n == 0 {
          n += 1
        }

        if errSub := checkOverload(overload[0:n-1], args, ctx); err != nil {
          context.AppendError(err, errSub)
        }
      }

      return nil, err
    }

    if proto, ok := v.interf.(Prototype); ok && proto.IsAbstract() && !allowAbstract {
      return nil, ctx.NewError("Error: can't construct abstract " + proto.Name())
    }
  }

  var interf Interface
  if v.fn != nil && args != nil {
    var err error
    interf, err = v.fn(args, ctx)
    if err != nil {
      return nil, err
    }
  } else if v.interf != nil {
    interf = v.interf
  }

  if interf == nil {
    return NewAny(ctx), nil
  } else {
    return NewInstance(interf, ctx), nil
  }
}

func (v *Class) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't call a class (hint: use new)")
}

func (v *Class) GetMember(key string, includePrivate bool,
  ctx context.Context) (Value, error) {
  proto := v.getPrototype()

  if proto == nil && v.interf == nil {
    return NewAny(ctx), nil
  } else if proto != nil {
    if res, err := proto.GetClassMember(key, includePrivate, ctx); err != nil {
      return nil, err
    } else if res != nil {
      return res, nil
    } else {
      return nil, ctx.NewError("Error: " + proto.Name() + "." + key + " not found")
    }
  } else {
    return nil, ctx.NewError("Error: static " + v.interf.Name() + "." + key + " not available because it is an interface")
  }
}

func (v *Class) SetMember(key string, includePrivate bool, arg Value,
  ctx context.Context) error {
  return ctx.NewError("Error: can't set static class members")
}

func IsClass(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  switch v_.(type) {
  case *Class:
    return true
  default:
    return false
  }
}
