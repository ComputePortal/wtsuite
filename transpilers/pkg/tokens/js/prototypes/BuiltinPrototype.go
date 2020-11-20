package prototypes

import (
  "../values"

  "../../context"
)

type BuiltinPrototype struct {
  name string
}

func newBuiltinPrototype(name string) BuiltinPrototype {
  return BuiltinPrototype{name}
}

func (p *BuiltinPrototype) Name() string {
  return p.name
}

func (p *BuiltinPrototype) Context() context.Context {
  return context.NewDummyContext()
}

/*func (p *BuiltinPrototype) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*BuiltinPrototype); ok {
    if p == other {
      return nil
    } else {
      return ctx.NewError("Error: expected " + p.Name() + ", got " + other_.Name())
    }
  } else if other, ok := other_.(values.Prototype); ok {
    if otherParent, err := other.GetParent(); err != nil {
      fmt.Println("1. ", other_.Name(), otherParent.Name())
      return err
    } else if otherParent != nil {
      fmt.Println("2. ", other_.Name(), otherParent.Name())
      if p.Check(otherParent, ctx) != nil {
        return ctx.NewError("Error: expected " + p.Name() + ", got " + other_.Name())
      } else {
        return nil
      }
    } else {
      fmt.Println("3. ", other_.Name(), otherParent.Name())
      return ctx.NewError("Error: expected " + p.Name() + ", got " + other_.Name())
    }
  } else {
    return ctx.NewError("Error: expected " + p.Name() + ", got " + other_.Name())
  }
}*/

func (p *BuiltinPrototype) GetInterfaces() ([]values.Interface, error) {
  return []values.Interface{}, nil
}

func (p *BuiltinPrototype) GetPrototypes() ([]values.Prototype, error) {
  return []values.Prototype{}, nil
}

func (p *BuiltinPrototype) IsUniversal() bool {
  return false
}

func (p *BuiltinPrototype) GetParent() (values.Prototype, error) {
  return nil, nil
}

func (p *BuiltinPrototype) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return nil, nil
}

func (p *BuiltinPrototype) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set " + p.Name() + "." + key)
}

func (p *BuiltinPrototype) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return nil, nil
}
