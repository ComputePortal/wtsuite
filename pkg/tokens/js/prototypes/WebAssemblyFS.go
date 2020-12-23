package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type WebAssemblyFS struct {
  AbstractBuiltinInterface
}

var webAssemblyFSInterface values.Interface = &WebAssemblyFS{newAbstractBuiltinInterface("WebAssemblyFS")}

func NewWebAssemblyFS(ctx context.Context) values.Value {
  return values.NewInstance(webAssemblyFSInterface, ctx)
}

// interfaces usually check in a different way, but this isn't really an interface
func (p *WebAssemblyFS) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebAssemblyFS); ok {
    return nil
  } else if other, ok := other_.(values.Prototype); ok {
    return checkInterfaceImplementation(p, other, 
    // list of getters
    []string{
      "close",
      "create",
      "exists",
      "open",
      "read",
      "seek",
      "size",
      "tell",
      "write",
    }, 
    // list of setters
    map[string]values.Value{}, 
    ctx)
  } else {
    return ctx.NewError("Error: not an WebAssemblyFS")
  }
}

func (p *WebAssemblyFS) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  switch key {
  case "close":
    return values.NewFunction([]values.Value{i, nil}, ctx), nil
  case "create", "open":
    return values.NewFunction([]values.Value{s, i}, ctx), nil
  case "exists":
    return values.NewFunction([]values.Value{s, b}, ctx), nil
  case "read":
    return values.NewFunction([]values.Value{i, i, NewUint8Array(ctx)}, ctx), nil
  case "seek":
    return values.NewFunction([]values.Value{i, i, nil}, ctx), nil
  case "size", "tell":
    return values.NewFunction([]values.Value{i, i}, ctx), nil
  case "write":
    return values.NewFunction([]values.Value{i, NewUint8Array(ctx), nil}, ctx), nil
  default:
    return nil, nil
  }
}
