package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Window struct {
  BuiltinPrototype
}

func NewWindowPrototype() values.Prototype {
  return &Window{newBuiltinPrototype("Window")}
}

func NewWindow(ctx context.Context) values.Value {
  return values.NewInstance(NewWindowPrototype(), ctx)
}

func (p *Window) GetParent() (values.Prototype, error) {
  return NewEventTargetPrototype(), nil
}

func (p *Window) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*Window); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func NewFetchFunction(ctx context.Context) values.Value {
  s := NewString(ctx)

  return values.NewFunction([]values.Value{s, NewPromise(NewResponse(ctx), ctx)}, ctx)
}

func NewSetTimeoutFunction(ctx context.Context) values.Value {
  fn := values.NewFunction([]values.Value{nil}, ctx)

  return values.NewOverloadedFunction([][]values.Value{
    []values.Value{fn, nil},
    []values.Value{fn, NewNumber(ctx), nil},
  }, ctx) 
}

func NewRequestIdleCallbackFunction(ctx context.Context) values.Value {
  fn := values.NewFunction([]values.Value{nil}, ctx)

  opt := NewConfigObject(map[string]values.Value{
    "timeout": NewNumber(ctx),
  }, ctx)

  return values.NewOverloadedFunction([][]values.Value{
    []values.Value{fn, nil},
    []values.Value{fn, opt, nil},
  }, ctx)
}

func (p *Window) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  f := NewNumber(ctx)
  s := NewString(ctx)

  switch key {
  case "atob", "btoa":
    return values.NewFunction([]values.Value{s, s}, ctx), nil
  case "blur", "close", "focus":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "crypto":
    return NewCrypto(ctx), nil
  case "devicePixelRatio", "innerHeight", "innerWidth", "scrollX", "scrollY":
    return f, nil
  case "fetch":
    return NewFetchFunction(ctx), nil
  case "getComputedStyle":
    return values.NewFunction([]values.Value{NewHTMLElement(ctx), NewCSSStyleDeclaration(ctx)}, ctx), nil
  case "indexedDB":
    return NewIDBFactory(ctx), nil
  case "localStorage", "sessionStorage":
    return NewStorage(ctx), nil
  case "location":
    return NewLocation(ctx), nil
  case "open":
    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{s, NewWindow(ctx)},
      []values.Value{s, s, NewWindow(ctx)},
    }, ctx), nil
  case "requestAnimationFrame":
    fn := values.NewFunction([]values.Value{f, nil}, ctx)

    return values.NewFunction([]values.Value{fn, nil}, ctx), nil
  case "requestIdleCallback":
    return NewRequestIdleCallbackFunction(ctx), nil
  case "scrollTo":
    return values.NewFunction([]values.Value{f, f, nil}, ctx), nil
  case "setInterval", "setTimeout":
    return NewSetTimeoutFunction(ctx), nil
  default:
    return nil, nil
  }
}

func (p *Window) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWindowPrototype(), ctx), nil
}
