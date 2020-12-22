package prototypes

import (
  "strings"

  "../values"

  "../../context"
)

type Tuple struct {
  content []values.Value // can be nil, or at least length 2

  BuiltinPrototype
}

func NewTuplePrototype(content []values.Value) values.Prototype {
  if content != nil {
    if len(content) < 2 {
      panic("should be at least length 2")
    }
  }

  return &Tuple{content, newBuiltinPrototype("Tuple")}
}

func NewTuple(content []values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewTuplePrototype(content), ctx)
}

func (p *Tuple) GetParent() (values.Prototype, error) {
  content := p.getCommonValue(p.Context())
  return NewArrayPrototype(content), nil
}

func (p *Tuple) IsUniversal() bool {
  if p.content == nil {
    return false
  } else {
    for _, content := range p.content {
      interf := values.GetInterface(content)
      if (interf != nil && (!interf.IsUniversal())) || interf == nil {
        return false
      } 
    }

    return true
  }
}

func (p *Tuple) getCommonValue(ctx context.Context) values.Value {
  if p.content == nil {
    return values.NewAny(ctx)
  } else {
    return values.CommonValue(p.content, ctx)
  }
}

func (p *Tuple) Name() string {
  var b strings.Builder

  b.WriteString("Tuple")

  if p.content != nil && len(p.content) > 0 {
    b.WriteString("<")

    for i, v := range p.content {
      b.WriteString(v.TypeName())

      if i < len(p.content) - 1 {
        b.WriteString(",")
      }
    }

    b.WriteString(">")
  }

  return b.String()
}

func (p *Tuple) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Tuple); ok {
    if p.content == nil || len(p.content) == 0 {
      return nil
    } else if other.content == nil || len(other.content) == 0 {
      return nil
    } else if len(p.content) != len(other.content) {
      return ctx.NewError("Error: tuples have different lengths")
    } else {
      for i, v := range p.content {
        otherV := other.content[i]
        if err := v.Check(otherV, ctx); err != nil {
          return err
        }
      }

      return nil
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Tuple) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)
  i := NewInt(ctx)
  common := p.getCommonValue(ctx)

  switch key {
  case ".getindex":
    return values.NewCustomFunction([]values.Value{i}, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
      if lit, ok := args[0].LiteralIntValue(); ok && p.content != nil {
        if lit < 0 || lit >= len(p.content) {
          errCtx := args[0].Context()
          return nil, errCtx.NewError("Error: out of tuple range")
        }

        return values.NewContextValue(p.content[lit], ctx_), nil
      } else {
        return common, nil
      }
    }, ctx), nil
  case ".setindex":
    return values.NewCustomFunction([]values.Value{i, a}, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
      if lit, ok := args[0].LiteralIntValue(); ok && p.content != nil {
        if lit < 0 || lit >= len(p.content) {
          errCtx := args[0].Context()
          return nil, errCtx.NewError("Error: out of tuple range")
        }

        checkVal := p.content[lit]
        return nil, checkVal.Check(args[1], args[1].Context())
      } else {
        return nil, common.Check(args[1], args[1].Context())
      }
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Tuple) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  a := values.NewAny(ctx)

  return values.NewCustomClass(
    [][]values.Value{
      []values.Value{a, a},
      []values.Value{a, a, a},
      []values.Value{a, a, a, a},
      []values.Value{a, a, a, a, a},
      []values.Value{a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a},
    }, func(args []values.Value, ctx_ context.Context) (values.Interface, error) {
      return NewTuplePrototype(args), nil
    }, ctx), nil
}
