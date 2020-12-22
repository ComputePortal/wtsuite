package prototypes

import (
  "strings"

  "../values"

  "../../context"
)

type Map struct {
  key values.Value // can be nil, then any
  item values.Value // can be nil, then any

  BuiltinPrototype
}

func NewMapPrototype(key values.Value, item values.Value) values.Prototype {
  return &Map{key, item, newBuiltinPrototype("Map")}
}

func NewMap(key values.Value, item values.Value, ctx context.Context) values.Value {
  return values.NewInstance(NewMapPrototype(key, item), ctx)
}

func (p *Map) getKeyValue(ctx context.Context) values.Value {
  if p.key == nil {
    return values.NewAny(ctx)
  } else {
    return values.NewContextValue(p.key, ctx)
  }
}

func (p *Map) getItemValue(ctx context.Context) values.Value {
  if p.item == nil {
    return values.NewAny(ctx)
  } else {
    return values.NewContextValue(p.item, ctx)
  }
}

func (p *Map) Name() string {
  var b strings.Builder

  b.WriteString("Map")

  if p.key != nil || p.item != nil {
    b.WriteString("<")

    if p.key == nil {
      b.WriteString("any")
    } else {
      b.WriteString(p.key.TypeName())
    }
    
    b.WriteString(",")

    if p.item == nil {
      b.WriteString("any")
    } else {
      b.WriteString(p.item.TypeName())
    }

    b.WriteString(">")
  }

  return b.String()
}

func (p *Map) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Map); ok {
    if p.key == nil && p.item == nil {
      return nil
    } else if other.key == nil && other.item == nil {
      return ctx.NewError("Error: expected " + p.Name() + ", got Map<any, any>")
    } else {
      if p.key != nil {
        if other.key == nil {
          return ctx.NewError("Error: expected " + p.Name() + ", got " + other.Name())
        } else if err := p.key.Check(other.key, ctx); err != nil {
          return err
        }
      }

      if p.item != nil {
        if other.item == nil {
          return ctx.NewError("Error: expected " + p.Name() + ", got " + other.Name())
        } else if err := p.item.Check(other.item, ctx); err != nil {
          return err
        } else {
          return nil
        }
      } else {
        return nil
      }
    }
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *Map) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  k := p.getKeyValue(ctx)
  item := p.getItemValue(ctx)

  switch key {
  case ".getof":
    return NewTuple([]values.Value{k, item}, ctx), nil
  case "clear":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "delete":
    return values.NewMethodLikeFunction([]values.Value{k, b}, ctx), nil
  case "get":
    return values.NewFunction([]values.Value{k, item}, ctx), nil
  case "set":
    return values.NewFunction([]values.Value{k, item, nil}, ctx), nil
  case "has":
    return values.NewFunction([]values.Value{k, b}, ctx), nil
  case "size":
    return i, nil
  default:
    return nil, nil
  }
}

func (p *Map) GetClassValue() (*values.Class, error) {
  ctx := p.Context()

  return values.NewClass(
    [][]values.Value{
      []values.Value{},
    }, NewMapPrototype(values.NewAny(ctx), values.NewAny(ctx)), ctx), nil
}
