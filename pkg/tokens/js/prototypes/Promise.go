package prototypes

import (
  "strings"

  "../values"

  "../../context"
)

type Promise struct {
  content values.Value // if nil, then any, also for reject (if not nil error for reject)

  isVoid bool

  BuiltinPrototype
}

func NewPromisePrototype(content values.Value) values.Prototype {
  return &Promise{content, false, newBuiltinPrototype("Promise")}
}

func NewVoidPromisePrototype() values.Prototype {
  return &Promise{nil, true, newBuiltinPrototype("Promise")}
}

func NewPromise(content values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewPromisePrototype(content), ctx)
}

func NewVoidPromise(ctx context.Context) values.Value {
  return values.NewInstance(NewVoidPromisePrototype(), ctx)
}

func (p *Promise) Name() string {
  var b strings.Builder

  b.WriteString("Promise")

  if p.isVoid {
    b.WriteString("<void>")
  } else if p.content != nil {
    b.WriteString("<")
    b.WriteString(p.content.TypeName())
    b.WriteString(">")
  }

  return b.String()
}

func IsPromise(v values.Value) bool {
  ctx := v.Context()

  checkVal := NewPromise(nil, ctx)

  return checkVal.Check(v, ctx) == nil
}

func (p *Promise) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Promise); ok {
    if p.content == nil {
      if p.isVoid && !other.isVoid {
        return ctx.NewError("Error: expected Promise<void>, got " + other.Name())
      }

      return nil
    } else if other.content == nil{
      return ctx.NewError("Error: expected " + p.Name() + ", got Promise<any>")
    } else if p.content.Check(other.content, ctx) != nil {
      return ctx.NewError("Error: expected " + p.Name() + ", got " + other.Name())
    } else {
      return nil
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Promise) getContent(ctx context.Context) values.Value {
  if p.content == nil {
    return values.NewAny(p.Context())
  } else {
    return values.NewContextValue(p.content, ctx)
  }
}

func (p *Promise) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)
  e := NewError(ctx)
  self := values.NewInstance(p, ctx)
  content := p.getContent(ctx)

  switch key {
  case "catch":
    if p.content == nil && !p.isVoid {
      return values.NewOverloadedMethodLikeFunction(
        [][]values.Value{
          []values.Value{values.NewFunction([]values.Value{nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{self}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, self}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, self}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, a, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, a, self}, ctx), self},
        }, ctx), nil
    } else {
      return values.NewOverloadedMethodLikeFunction(
        [][]values.Value{
          []values.Value{values.NewFunction([]values.Value{e, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{e, self}, ctx), self},
        }, ctx), nil
    }
  case "then":
    if p.isVoid {
      return values.NewOverloadedMethodLikeFunction(
        [][]values.Value{
          []values.Value{values.NewFunction([]values.Value{nil}, ctx), self},
        }, ctx), nil
    } else if p.content == nil {
      return values.NewOverloadedMethodLikeFunction(
        [][]values.Value{
          []values.Value{values.NewFunction([]values.Value{nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{self}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, self}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, self}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, a, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{a, a, a, self}, ctx), self},
        }, ctx), nil
    } else {
      return values.NewOverloadedMethodLikeFunction(
        [][]values.Value{
          []values.Value{values.NewFunction([]values.Value{content, nil}, ctx), self},
          []values.Value{values.NewFunction([]values.Value{content, self}, ctx), self},
        }, ctx), nil
    }
  case ".resolve":
    if p.isVoid {
      return values.NewFunction([]values.Value{nil}, ctx), nil
    } else {
      return values.NewFunction([]values.Value{content}, ctx), nil
    }
  default:
    return nil, nil
  }
}

func GetPromiseContent(v values.Value) (values.Value, error) {
  resolveFn, err := v.GetMember(".resolve", false, v.Context())
  if err != nil {
    return nil, err
  }

  return resolveFn.EvalFunction([]values.Value{}, false, v.Context())
}

func (p *Promise) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "all":
    return values.NewOverloadedCustomFunction([][]values.Value{
      []values.Value{NewArray(NewPromise(values.NewAny(ctx), ctx), ctx)},
      []values.Value{NewArray(NewVoidPromise(ctx), ctx)},
    }, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
      prom, err := args[0].GetMember(".getof", false, ctx)
      if err != nil {
        return nil, err
      }

      resFn, err := prom.GetMember(".resolve", false, ctx)
      if err != nil {
        return nil, err
      }

      res, err := resFn.EvalFunction([]values.Value{}, false, ctx)
      if err != nil {
        return nil, err
      }

      if res == nil {
        return NewVoidPromise(ctx_), nil
      } else {
        return NewPromise(NewArray(res, ctx_), ctx_), nil
      }
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Promise) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  a := values.NewAny(ctx)
  e := NewError(ctx)

  return values.NewCustomClass(
    [][]values.Value{
      []values.Value{values.NewFunction([]values.Value{
        values.NewFunction([]values.Value{a, nil}, ctx),
        values.NewFunction([]values.Value{e, nil}, ctx),
      nil}, ctx)},
      []values.Value{values.NewFunction([]values.Value{
        values.NewFunction([]values.Value{nil}, ctx),
        values.NewFunction([]values.Value{e, nil}, ctx),
      nil}, ctx)},
    }, func(args []values.Value, ctx_ context.Context) (values.Interface, error) {
      if args == nil {
        return NewPromisePrototype(values.NewAny(ctx)), nil
      } else {
        val, err := args[0].GetMember(".arg0", false, ctx)
        if err != nil {
          panic(err)
        }

        val, err = val.GetMember(".arg0", false, ctx)
        if err == nil {
          if val == nil {
            return NewVoidPromisePrototype(), nil
          } else {
            return NewPromisePrototype(val), nil
          }
        } 

        return NewPromisePrototype(values.NewAny(ctx)), nil
      }
    }, ctx), nil
}
