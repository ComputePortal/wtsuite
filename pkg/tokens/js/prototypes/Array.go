package prototypes

import (
  "strings"

  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

// instance of Array are values.Prototype, not values.Value!
type Array struct {
  content values.Value // if nil, then any
  BuiltinPrototype
}

func NewArrayPrototype(content values.Value) values.Prototype {
  return &Array{content, newBuiltinPrototype("Array")}
}

func NewArray(content values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewArrayPrototype(content), ctx)
}

func IsArray(v values.Value) bool {
  ctx := context.NewDummyContext()

  arrayCheck := NewArray(nil, ctx)

  return arrayCheck.Check(v, ctx) == nil
}

func (p *Array) getContent(ctx context.Context) values.Value {
  if p.content == nil {
    return values.NewAny(ctx)
  } else {
    return values.NewContextValue(p.content, ctx)
  }
}

func (p *Array) Name() string {
  var b strings.Builder

  b.WriteString("Array")

  if p.content != nil {
    b.WriteString("<")
    b.WriteString(p.content.TypeName())
    b.WriteString(">")
  }

  return b.String()
}

func (p *Array) IsUniversal() bool {
  if p.content == nil {
    return false
  } else if interf := values.GetInterface(p.content); interf != nil {
    return interf.IsUniversal()
  } else {
    return false
  }
}

func (p *Array) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Array); ok {
    if p.content == nil {
      return nil
    } else if other.content == nil{
      return ctx.NewError("Error: expected Array<" + p.content.TypeName() + ">, got Array<any>")
    } else if p.content.Check(other.content, ctx) != nil {
      return ctx.NewError("Error: expected Array<" + p.content.TypeName() + ">, got Array<" + other.content.TypeName() + ">")
    } else {
      return nil
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Array) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  // some commonly used values
  a := values.NewAny(ctx)
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  f := NewNumber(ctx)
  s := NewString(ctx)

  content := p.getContent(ctx)
  self := NewArray(content, ctx)

  switch key {
  case ".getindex":
    return values.NewFunction([]values.Value{i, content}, ctx), nil
  case ".setindex":
    // TODO: use this, instead of basing on .getindex
    return values.NewFunction([]values.Value{i, content, nil}, ctx), nil
  case ".getin":
    return nil, ctx.NewError("Error: can't iterate over Array using 'in' (hint: use regular 'for', or 'for of' if your are only interested in the items)")
  case ".getof":
    return content, nil
  case "concat":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{self, self},
        []values.Value{self, self, self},
        []values.Value{self, self, self, self},
        []values.Value{self, self, self, self, self},
        []values.Value{self, self, self, self, self, self}, // 5 should be enough
      }, ctx), nil
  case "copyWithin":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{i, nil},
        []values.Value{i, i, nil},
        []values.Value{i, i, i, nil},
      }, ctx), nil
  case "any", "every", "some":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{values.NewFunction([]values.Value{content, b}, ctx), b}, 
        []values.Value{values.NewFunction([]values.Value{content, i, b}, ctx), b}, 
        []values.Value{values.NewFunction([]values.Value{content, i, self, b}, ctx), b}, 
      }, ctx), nil
  case "fill":
    return values.NewOverloadedMethodLikeFunction(
      [][]values.Value{
        []values.Value{content, self},
        []values.Value{content, i, self},
        []values.Value{content, i, i, self},
      }, ctx), nil
  case "find":
    return nil, ctx.NewError("Error: use findIndex instead (find returns undefined, which is not (yet) supported)")
  case "findIndex":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{values.NewFunction([]values.Value{content, b}, ctx), i}, 
        []values.Value{values.NewFunction([]values.Value{content, i, b}, ctx), i}, 
        []values.Value{values.NewFunction([]values.Value{content, i, self, b}, ctx), i}, 
      }, ctx), nil
  case "filter":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{values.NewFunction([]values.Value{content, b}, ctx), self},
        []values.Value{values.NewFunction([]values.Value{content, i, b}, ctx), self},
        []values.Value{values.NewFunction([]values.Value{content, i, self, b}, ctx), self},
      }, ctx), nil
  case "forEach":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{values.NewFunction([]values.Value{content, nil}, ctx), nil},
        []values.Value{values.NewFunction([]values.Value{content, i, nil}, ctx), nil},
        []values.Value{values.NewFunction([]values.Value{content, i, self, nil}, ctx), nil},
      }, ctx), nil
  case "indexOf", "lastIndexOf": 
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{content, i},
        []values.Value{content, i, i},
      }, ctx), nil
  case "join":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{s},
        []values.Value{s, s},
      }, ctx), nil
  case "length":
    return i, nil
  case "map":
    return values.NewOverloadedCustomFunction(
      [][]values.Value{
        []values.Value{values.NewFunction([]values.Value{content, a}, ctx)},
        []values.Value{values.NewFunction([]values.Value{content, i, a}, ctx)},
        []values.Value{values.NewFunction([]values.Value{content, i, self, a}, ctx)},
      }, func(args []values.Value, preferMethod bool, ctx_ context.Context) (values.Value, error) {
        ret, err := args[0].GetMember(".return", false, ctx_)
        if err != nil {
          panic("expected function")
        }

        if ret == nil {
          panic("expected return value")
        }
        
        return NewArray(ret, ctx), nil
      }, ctx), nil
  case "pop", "shift":
    return values.NewFunction([]values.Value{content}, ctx), nil
  case "push", "unshift":
    return values.NewMethodLikeFunction([]values.Value{content, self}, ctx), nil
  case "reduce", "reduceRight":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{values.NewFunction([]values.Value{content, content, content}, ctx), content},
        []values.Value{values.NewFunction([]values.Value{content, content, i, content}, ctx), content},
        []values.Value{values.NewFunction([]values.Value{content, content, i, self, content}, ctx), content},
        []values.Value{values.NewFunction([]values.Value{content, content, content}, ctx), content, content},
        []values.Value{values.NewFunction([]values.Value{content, content, i, content}, ctx), content, content},
        []values.Value{values.NewFunction([]values.Value{content, content, i, self, content}, ctx), content, content},
      }, ctx), nil
  case "reverse":
    return values.NewMethodLikeFunction([]values.Value{self}, ctx), nil
  case "slice":
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{self},
        []values.Value{i, self},
        []values.Value{i, i, self},
      }, ctx), nil
  case "sort":
    return values.NewOverloadedMethodLikeFunction(
      [][]values.Value{
        []values.Value{self},
        []values.Value{values.NewFunction([]values.Value{content, content, i}, ctx), self},
        []values.Value{values.NewFunction([]values.Value{content, content, f}, ctx), self},
      }, ctx), nil
  case "splice":
    return values.NewOverloadedMethodLikeFunction(
      [][]values.Value{
        []values.Value{i, self},
        []values.Value{i, i, self},
        []values.Value{i, i, content, self},
        []values.Value{i, i, content, content, self},
        []values.Value{i, i, content, content, content, self}, // 3 should be enough
      }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Array) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set Array." + key)
}

func (p *Array) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  a := values.NewAny(ctx)
  b := NewBoolean(ctx)

  switch key {
  case "isArray":
    return values.NewFunction([]values.Value{a, b}, ctx), nil
  case "from":
    // TODO: complete this interface
    return values.NewOverloadedFunction(
      [][]values.Value{
        []values.Value{NewSet(nil, ctx), NewArray(nil, ctx)},
        []values.Value{NewArray(nil, ctx), NewArray(nil, ctx)},
        []values.Value{NewMap(nil, nil, ctx), NewArray(nil, ctx)},
      }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *Array) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  a := values.NewAny(ctx)
  i := NewInt(ctx)

  return values.NewClass(
    [][]values.Value{
      []values.Value{},
      []values.Value{i},
      []values.Value{a, a},
      []values.Value{a, a, a},
      []values.Value{a, a, a, a},
      []values.Value{a, a, a, a, a},
      []values.Value{a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
      []values.Value{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
    }, NewArrayPrototype(values.NewAny(ctx)), ctx), nil
}
