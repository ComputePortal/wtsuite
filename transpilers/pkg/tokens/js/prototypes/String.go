package prototypes

import (
  "../values"

  "../../context"
)

type String struct {
  BuiltinPrototype
}

func NewStringPrototype() values.Prototype {
  return &String{newBuiltinPrototype("String")}
}

func NewString(ctx context.Context) values.Value {
  return values.NewInstance(NewStringPrototype(), ctx)
}

func NewLiteralString(v string, ctx context.Context) values.Value {
  return values.NewLiteralStringInstance(NewStringPrototype(), v, ctx)
}

func IsString(v values.Value) bool {
  ctx := context.NewDummyContext()

  stringCheck := NewString(ctx)

  return stringCheck.Check(v, ctx) == nil
}

func IsStringable(v values.Value) bool {
  return IsString(v) || IsNumber(v) || IsBoolean(v)
}

func (p *String) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*String); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *String) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  re := NewRegExp(ctx)
  s := NewString(ctx)
  ss := NewArray(s, ctx)

  switch key {
  case ".getindex", "charAt":
    return values.NewFunction([]values.Value{i, s}, ctx), nil
  case ".getof":
    return s, nil
  case "charCodeAt", "codePointAt":
    return values.NewFunction([]values.Value{i, i}, ctx), nil
  case "concat":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, s},
      []values.Value{s, s, s},
      []values.Value{s, s, s, s},
      []values.Value{s, s, s, s, s},
    }, ctx), nil
  case "endsWith", "includes", "startsWith":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, b},
      []values.Value{s, i, b},
    }, ctx), nil
  case "indexOf", "lastIndexOf":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, i},
      []values.Value{s, i, i},
    }, ctx), nil
  case "length":
    return i, nil
  case "localeCompare":
    opt := NewLocaleOptions(ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, i},
      []values.Value{s, s, i},
      []values.Value{s, s, opt, i},
    }, ctx), nil
  case "match":
    return values.NewFunction([]values.Value{NewRegExp(ctx), ss}, ctx), nil
  case "normalize":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{s, s},
    }, ctx), nil
  case "padEnd", "padStart":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, s},
      []values.Value{i, s, s},
    }, ctx), nil
  case "repeat":
    return values.NewFunction([]values.Value{i, s}, ctx), nil
  case "replace":
    // non regexp functions
    fn1 := values.NewFunction([]values.Value{s, s}, ctx)
    fn2 := values.NewFunction([]values.Value{s, i, s}, ctx)
    fn3 := values.NewFunction([]values.Value{s, i, s, s}, ctx)

    fn4 := values.NewFunction([]values.Value{s, s, i, s, s}, ctx)
    fn5 := values.NewFunction([]values.Value{s, s, s, i, s, s}, ctx)
    fn6 := values.NewFunction([]values.Value{s, s, s, s, i, s, s}, ctx) // 3 capture groups should be enough
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, s, s},
      []values.Value{s, fn1, s},
      []values.Value{s, fn2, s},
      []values.Value{s, fn3, s},
      []values.Value{re, s, s},
      []values.Value{re, fn1, s},
      []values.Value{re, fn2, s},
      []values.Value{re, fn3, s},
      []values.Value{re, fn4, s},
      []values.Value{re, fn5, s},
      []values.Value{re, fn6, s},
    }, ctx), nil
  case "search":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, i},
      []values.Value{re, i},
    }, ctx), nil
  case "slice", "substring":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, s},
      []values.Value{i, i, s},
    }, ctx), nil
  case "split":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, ss},
      []values.Value{s, i, ss},
      []values.Value{re, ss},
      []values.Value{re, i, ss},
    }, ctx), nil
  case "toLocaleLowerCase", "toLocaleUpperCase":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s},
      []values.Value{s, s},
      []values.Value{ss, s},
    }, ctx), nil
  case "toLowerCase", "toUpperCase", "trim", "trimLeft", "trimRight":
    return values.NewFunction([]values.Value{s}, ctx), nil
  default:
    return nil, nil
  }
}

func (p *String) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "fromCharCode", "fromCodePoint":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{i, s},
      []values.Value{i, i, s},
      []values.Value{i, i, i, s},
      []values.Value{i, i, i, i, s},
      []values.Value{i, i, i, i, i, s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *String) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  a := values.NewAny(ctx)

  return values.NewClass([][]values.Value{
    []values.Value{a},
  }, NewStringPrototype(), ctx), nil
}

func (p *String) IsUniversal() bool {
  return true
}
