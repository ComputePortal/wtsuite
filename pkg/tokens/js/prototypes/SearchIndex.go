package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SearchIndex struct {
  BuiltinPrototype
}

func NewSearchIndexPrototype() values.Prototype {
  return &SearchIndex{newBuiltinPrototype("SearchIndex")}
}

func NewSearchIndex(ctx context.Context) values.Value {
  return values.NewInstance(NewSearchIndexPrototype(), ctx)
}

func (p *SearchIndex) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*SearchIndex); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *SearchIndex) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)
  ss := NewArray(s, ctx)

  switch key {
  case "onready":
    return nil, ctx.NewError("Error: only a setter")
  case "ignore":
    return values.NewFunction([]values.Value{s, b}, ctx), nil
  case "page":
    return values.NewFunction([]values.Value{i, NewObject(map[string]values.Value{
      "url": s,
      "title": s,
      "content": ss,
    }, ctx)}, ctx), nil
  case "match", "matchPrefix", "matchSuffix", "matchSubstring":
    return values.NewFunction([]values.Value{s, NewSet(i, ctx)}, ctx), nil
  case "fuzzy", "fuzzyPrefix", "fuzzySuffix", "fuzzySubstring":
    return values.NewFunction([]values.Value{s, i, NewArray(NewSet(i, ctx), ctx)}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *SearchIndex) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  switch key {
  case "onready":
    callback := values.NewFunction([]values.Value{nil}, ctx)
    return callback.Check(arg, ctx)
  default:
    return ctx.NewError("Error: SearchIndex." + key + " not setable")
  }
}

func (p *SearchIndex) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  s := NewString(ctx)

  opt := NewConfigObject(map[string]values.Value{
  }, ctx)

  return values.NewClass([][]values.Value{
    []values.Value{s},
    []values.Value{s, opt},
  }, NewSearchIndexPrototype(), ctx), nil
}
