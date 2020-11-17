package prototypes

import (
  "../values"

  "../../context"
)

type Object struct {
  isConfig bool

  common values.Value

  members map[string]values.Value // can be nil

  BuiltinPrototype
}

func NewObjectPrototype(members map[string]values.Value) values.Prototype {
  obj := &Object{false, nil, members, newBuiltinPrototype("Object")}

  obj.common = values.CommonValue(obj.getValues())

  return obj
}

func NewMapLikeObjectPrototype(common values.Value) values.Prototype {
  return &Object{false, common, nil, newBuiltinPrototype("Object")}
}

func NewObject(members map[string]values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewObjectPrototype(members), ctx)
}

func NewMapLikeObject(common values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewMapLikeObjectPrototype(common), ctx)
}

func NewConfigObject(members map[string]values.Value, ctx context.Context) values.Value {
  if members == nil {
    panic("can't be nil")
  }

  proto := &Object{true, nil, members, newBuiltinPrototype("Object")}

  proto.common = values.CommonValue(proto.getValues())

  return values.NewInstance(proto, ctx)
}

func (p *Object) IsUniversal() bool {
  if p.common == nil {
    if p.members == nil {
      return false
    } else {
      for _, v := range p.members {
        vProto := values.GetPrototype(v)
        if vProto == nil {
          return false
        } else if !vProto.IsUniversal() {
          return false
        }
      }

      return true
    }
  } else if proto := values.GetPrototype(p.common); proto != nil {
    return proto.IsUniversal()
  } else {
    return false
  }
}

func (p *Object) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Object); ok {
    if p.members == nil {
      if p.common == nil {
        return nil
      } else {
        if err := p.common.Check(other.common, ctx); err != nil {
          return err
        }

        return nil
      }
    } else if other.members == nil { // XXX: should we be more permissive
      return ctx.NewError("Error: expected Object with typed content")
    } else if p.isConfig && other.members != nil {
      for k, v := range other.members {
        if thisV, ok := p.members[k]; ok {
          if err := thisV.Check(v, ctx); err != nil {
            return ctx.NewError("Error: option " + k + " has invalid type (expected " + thisV.TypeName() + ", got " + v.TypeName() + ")")
          }
        } else {
          return ctx.NewError("Error: unrecognized option " + k)
        }
      }

      return nil
    } else {
      for k, v := range p.members {
        if otherV, ok := other.members[k]; !ok {
          return ctx.NewError("Error: missing Object." + k)
        } else {
          if err := v.Check(otherV, ctx); err != nil {
            return err
          }
        }
      }

      return nil
    }
  } else {
    return ctx.NewError("Error: expected Object, got " + other_.Name())
  }
}

func (p *Object) getMember(k string, ctx context.Context) (values.Value, error) {
  if p.members != nil {
    if v, ok := p.members[k]; ok {
      return v, nil
    } else {
      return nil, ctx.NewError("Error: Object." + k + " not found")
    }
  }

  common := p.getCommonValue()

  return values.NewContextValue(common, ctx), nil
}

func (p *Object) getValues() []values.Value {
  if p.members != nil {
    vs := make([]values.Value, 0)

    for _, v := range p.members {
      vs = append(vs, v)
    }

    return vs
  } else {
    return []values.Value{values.NewAny(p.Context())}
  }
}

func (p *Object) getCommonValue() values.Value {
  return p.common
}

func (p *Object) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)
  s := NewString(ctx)
  common := p.getCommonValue()

  switch key {
  case ".getindex":
    return values.NewCustomFunction([]values.Value{s}, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
      if k, ok := args[0].LiteralStringValue(); ok {
        return p.getMember(k, ctx_)
      } else {
        return values.NewContextValue(common, ctx_), nil
      }
    }, ctx), nil
  case ".setindex":
    return values.NewCustomFunction([]values.Value{s, a}, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
      if k, ok := args[0].LiteralStringValue(); ok {
        m, err := p.getMember(k, ctx_)
        if err != nil {
          return nil, err
        }

        if err := m.Check(args[1], ctx_); err != nil {
          return nil, err
        }
      } else {
        if err := common.Check(args[1], ctx_); err != nil {
          return nil, err
        }
      }

      return nil, nil
    }, ctx), nil
  case ".getof":
    return nil, ctx.NewError("Error: can't loop over Object using 'for of' (hint: use 'for in')")
  case ".getin":
    return s, nil
  default:
    if p.members == nil {
      return nil, nil
    } else {
      v, ok := p.members[key]
      if !ok {
        return nil, nil
      }

      return values.NewContextValue(v, ctx), nil
    }
  }
}

func (p *Object) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  if p.members != nil {
    checkVal, ok := p.members[key]
    if ok {
      return checkVal.Check(arg, ctx)
    } else {
      return ctx.NewError("Error: member " + key + " not found")
    }
  } else {
    return ctx.NewError("Error: object has unknown content (hint: use [...] instead)")
  }
}

func (p *Object) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  o := NewObject(nil, ctx)

  switch key {
  case "assign":
    return values.NewCustomFunction(
      []values.Value{
        o, o,
      }, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
        proto1_ := values.GetPrototype(args[0])
        proto2_ := values.GetPrototype(args[1])

        proto1, ok1 := proto1_.(*Object)
        proto2, ok2 := proto2_.(*Object)

        if ok1 && ok2 {
          if proto1.members != nil && proto2.members != nil {
            members := make(map[string]values.Value)
            for k, m := range proto1.members {
              members[k] = m
            }
            for k, m := range proto2.members {
              members[k] = m
            }

            return NewObject(members, ctx_), nil
          } else {
            common := values.CommonValue([]values.Value{proto1.common, proto2.common})

            return NewMapLikeObject(common, ctx_), nil
          }
        } else {
          return values.NewContextValue(o, ctx_), nil
        }
      }, ctx), nil
  case "keys":
    return values.NewFunction(
      []values.Value{
        NewObject(nil, ctx), NewArray(NewString(ctx), ctx),
      }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Object) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewClass([][]values.Value{
    []values.Value{},
  }, NewMapLikeObjectPrototype(values.NewAll(ctx)), ctx), nil
}
