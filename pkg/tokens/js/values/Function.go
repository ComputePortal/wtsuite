package values

import (
  "fmt"
  "strings"
  "strconv"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Function struct {
  methodLike bool 

  args [][]Value // last value of each list is the return value

  fn func(args []Value, preferMethod bool, ctx context.Context) (Value, error)

  ValueData
}

func NewFunction(argsAndRet []Value, ctx context.Context) *Function {
  return NewOverloadedFunction([][]Value{argsAndRet}, ctx)
}

func NewMethodLikeFunction(argsAndRet []Value, ctx context.Context) *Function {
  return NewOverloadedMethodLikeFunction([][]Value{argsAndRet}, ctx)
}

func NewOverloadedFunction(argsAndRet [][]Value, ctx context.Context) *Function {
  return &Function{false, argsAndRet, nil, ValueData{ctx}}
}

func NewOverloadedMethodLikeFunction(argsAndRet [][]Value, ctx context.Context) *Function {
  return &Function{true, argsAndRet, nil, ValueData{ctx}}
}

// args dont contain return value in case fn is defined
func NewCustomFunction(args []Value, fn func(args []Value, preferMethod bool, ctx_ context.Context) (Value, error), ctx context.Context) *Function {
  return NewOverloadedCustomFunction([][]Value{args}, fn, ctx)
}

func NewOverloadedCustomFunction(args [][]Value, fn func(args []Value, preferMethod bool, ctx_ context.Context) (Value, error), ctx context.Context) *Function {
  return &Function{false, args, fn, ValueData{ctx}}
}

func (v *Function) TypeName() string {
  var b strings.Builder

  b.WriteString("function")

  if v.args != nil {
    b.WriteString("<")

    if len(v.args) == 1 {
      for i, arg := range v.args[0] {

        if i < len(v.args[0]) - 1 {
          b.WriteString(arg.TypeName())
          b.WriteString(",")
        } else if arg == nil {
          b.WriteString("void")
        } else {
          b.WriteString(arg.TypeName())
        }
      }
    } else {
      b.WriteString(fmt.Sprintf("%d overloads", len(v.args)))
    }

    b.WriteString(">")
  }

  return b.String()
}

func (v *Function) IsVoid() bool {
  for _, args := range v.args {
    ret := args[len(args)-1]
    if ret != nil {
      return false
    }
  }

  return true
}

// get args, but ignore the return value
func (v *Function) GetArgs() [][]Value {
  args := make([][]Value, len(v.args))

  for i, overload := range v.args {
    args[i] = overload[0:len(overload)-1] // cut off of the return value
  }

  return args
}

func (v *Function) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if IsAny(other_) {
    return nil
  } else if other, ok := other_.(*Function); ok {
    if v.args != nil {
      if other.args == nil {
        errCtx := ctx
        return errCtx.NewError("Error: other function is generic")
      }

      if err := checkAllOverloads(v.args, other.args, ctx); err != nil {
        return err
      }
    }

    return nil
  } else {
    return ctx.NewError("Error: not a function")
  }
}

func (v *Function) EvalConstructor(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: function cannot be constructed (hint: remove new)")
}

// in case v.fn isn't defined
func (v *Function) evalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  // args don't include return value
  if preferMethod {
    // first check overloads with void return value
    for _, overload := range v.args {
      n := len(overload)
      if overload[n-1] == nil {
        if err := checkOverload(overload[0:n-1], args, ctx); err == nil {
          return nil, nil
        } else if len(v.args) == 1 {
          return nil, err
        }
      }
    }
  } else {
    for _, overload := range v.args {
      n := len(overload)
      if overload[n-1] != nil {
        if err := checkOverload(overload[0:n-1], args, ctx); err == nil {
          return NewContextValue(overload[n-1], ctx), nil
        } else if len(v.args) == 1 {
          return nil, err
        }
      }
    }
  }

  for _, overload := range v.args {
    n := len(overload)

    if err := checkOverload(overload[0:n-1], args, ctx); err == nil {
      if overload[n-1] != nil {
        return NewContextValue(overload[n-1], ctx), nil
      } else {
        return nil, nil
      }
    } else if len(v.args) == 1 {
      return nil, err
    }
  }

  err := ctx.NewError("Error: function arg types differ")
  
  // add info for each overload
  for _, overload := range v.args {
    n := len(overload)

    if errSub := checkOverload(overload[0:n-1], args, ctx); err != nil {
      context.AppendError(err, errSub)
    }
  }

  return nil, err
}

func (v *Function) EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) {
  if v.args == nil {
    if v.fn != nil {
      return v.fn(args, preferMethod, ctx)
    } else if preferMethod {
      return nil, nil
    } else {
      return NewAny(ctx), nil
    }
  }

  var ret Value = nil
  if v.fn == nil {
    var err error
    ret, err = v.evalFunction(args, preferMethod, ctx)
    if err != nil {
      return nil, err
    }
  } else {
    _, err := checkAnyOverload(v.args, args, ctx)
    if err != nil {
      return nil, err
    }

    ret, err = v.fn(args, preferMethod, ctx)
    if err != nil {
      return nil, err
    }
  }

  if ret != nil && preferMethod && v.methodLike {
    return nil, nil
  } else {
    return ret, nil
  }
}

func (v *Function) GetMember(key string, includePrivate bool, ctx context.Context) (Value, error) {
  if strings.HasPrefix(key, ".arg") {
    i, err := strconv.Atoi(key[4:])
    if err != nil {
      panic(err)
    }

    for _, overload := range v.args {
      if i < len(overload) {
        return overload[i], nil
      }
    }

    return nil, ctx.NewError("Error: " + key + " not found")
  } else if key == ".return" {
    overload := v.args[0]
    n := len(overload)
    return overload[n-1], nil
  }

  return nil, ctx.NewError("Error: can't get member of function")
}

func (v *Function) SetMember(key string, includePrivate bool, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set member of function")
}
