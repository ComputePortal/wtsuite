package values

import (
  "fmt"
  "strings"

	"../../context"
)

type Class struct {
  args [][]Value // various constructor overloads, can be nil in generic case

  fn func(args []Value, ctx_ context.Context) (Prototype, error)

  proto Prototype

	ValueData
}

func NewClass(args [][]Value, proto Prototype, ctx context.Context) *Class {
	return &Class{args, nil, proto, ValueData{ctx}}
}

func NewCustomClass(args [][]Value, fn func(args []Value, ctx_ context.Context) (Prototype, error), ctx context.Context) *Class {
  return &Class{args, fn, nil, ValueData{ctx}}
}

func NewUnconstructableClass(proto Prototype, ctx context.Context) *Class {
  return &Class{[][]Value{[]Value{}}, func(args []Value, ctx_ context.Context) (Prototype, error) {
    if args == nil {
      return proto, nil
    } else {
      return nil, ctx_.NewError("Error: doesn't have a constructor")
    }
  }, nil, ValueData{ctx}}
}

func (v *Class) GetConstructorArgs() [][]Value {
  return v.args
}

func (v *Class) getPrototype() Prototype {
  if v.proto == nil {
    if v.fn != nil {
      proto, err := v.fn(nil, context.NewDummyContext())
      if err != nil {
        panic(err)
      }

      return proto
    } else {
      return nil
    }
  } else {
    return v.proto
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

    proto := v.getPrototype()
    b.WriteString(proto.Name())

    b.WriteString(">")
  }

  return b.String()
}

// maybe it is a little silly that GetType always needs to be called
func (v *Class) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAll(other_) {
    return nil
  } else if other, ok := other_.(*Class); ok {
    if err := checkAllOverloads(v.args, other.args, ctx); err != nil {
      return err
    }

    proto := v.getPrototype()
    otherProto := other.getPrototype()
    if proto == nil {
      return nil
    } else if otherProto == nil {
      return ctx.NewError("Error: unspecified prototype")
    } else {
      return proto.Check(otherProto, ctx)
    }
  } else {
    return ctx.NewError("Error: not a class")
  }
}

func (v *Class) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  if _, err := checkAnyOverload(v.args, args, ctx); err != nil {
    return nil, err
  }

  var proto Prototype
  if v.fn != nil {
    var err error
    proto, err = v.fn(args, ctx)
    if err != nil {
      return nil, err
    }
  } else if v.proto != nil {
    proto = v.proto
  }

  if proto == nil {
    return NewAny(ctx), nil
  } else {
    return NewInstance(proto, ctx), nil
  }
}

func (v *Class) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't call a class (hint: use new)")
}

func (v *Class) GetMember(key string, includePrivate bool,
  ctx context.Context) (Value, error) {
  proto := v.getPrototype()

  if proto == nil {
    return NewAny(ctx), nil
  } else {
    return proto.GetClassMember(key, includePrivate, ctx)
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
