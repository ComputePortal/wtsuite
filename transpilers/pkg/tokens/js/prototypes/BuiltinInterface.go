package prototypes

import (
  "../values"

  "../../context"
)

type BuiltinInterface interface {
  values.Interface

  hasImplementation(values.Prototype) bool

  appendImplementation(values.Prototype)
}

type AbstractBuiltinInterface struct {
  name string

  // might need mutex in case of parallel eval
  implementations []values.Prototype
}

func newAbstractBuiltinInterface(name string) AbstractBuiltinInterface {
  return AbstractBuiltinInterface{name, make([]values.Prototype, 0)}
}

func (p *AbstractBuiltinInterface) Name() string {
  return p.name
}

func (p *AbstractBuiltinInterface) Context() context.Context {
  return context.NewDummyContext()
}

func (p *AbstractBuiltinInterface) IsUniversal() bool {
  for _, proto := range p.implementations {
    if !proto.IsUniversal() {
      return false
    }
  }

  return true
}

func (p *AbstractBuiltinInterface) GetInterfaces() ([]values.Interface, error) {
  return []values.Interface{}, nil
}

func (p *AbstractBuiltinInterface) GetPrototypes() ([]values.Prototype, error) {
  return p.implementations, nil
}

func (p *AbstractBuiltinInterface) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  return nil, nil
}

func (p *AbstractBuiltinInterface) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  return ctx.NewError("Error: " + p.Name() + "." + key + " not setable")
}

func (p *AbstractBuiltinInterface) hasImplementation(proto values.Prototype) bool {
  for _, ips := range p.implementations {
    if ips == proto {
      return true
    }
  }

  return false
}

func (p *AbstractBuiltinInterface) appendImplementation(proto values.Prototype) {
  p.implementations = append(p.implementations, proto)
}

func checkInterfaceImplementation(interf BuiltinInterface, proto values.Prototype, getables[]string, setables map[string]values.Value, ctx context.Context) error {
  if interf.hasImplementation(proto) {
    return nil
  }

  for _, key := range getables {
    v, err := interf.GetInstanceMember(key, false, ctx)
    if err != nil {
      panic(err)
    }

    vProto, err := proto.GetInstanceMember(key, false, ctx)
    if err != nil {
      return err
    }

    if err := v.Check(vProto, ctx); err != nil {
      return err
    }
  }

  for key, v := range setables {
    if err := proto.SetInstanceMember(key, false, v, ctx); err != nil {
      return err
    }
  }

  interf.appendImplementation(proto)

  return nil
}
