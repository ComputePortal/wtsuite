package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type HTMLCanvasElement struct {
  BuiltinPrototype
}

func NewHTMLCanvasElementPrototype() values.Prototype {
  return &HTMLCanvasElement{newBuiltinPrototype("HTMLCanvasElement")}
}

func NewHTMLCanvasElement(ctx context.Context) values.Value {
  return values.NewInstance(NewHTMLCanvasElementPrototype(), ctx)
}

func (p *HTMLCanvasElement) GetParent() (values.Prototype, error) {
  return NewHTMLElementPrototype(), nil
}

func (p *HTMLCanvasElement) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*HTMLCanvasElement); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *HTMLCanvasElement) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  f := NewNumber(ctx)
  s := NewString(ctx)

  switch key {
  case "getContext":
    o2d := NewConfigObject(map[string]values.Value{
      "alpha": b,
      "desynchronized": b,
    }, ctx)

    ogl := NewConfigObject(map[string]values.Value{
      "alpha": b,
      "desynchronized": b,
      "antialias": b,
      "depth": b,
      "failIfMajorPerformanceCaveat": b,
      "powerPreference": s,
      "premultipliedAlpha": b,
      "preserveDrawingBuffer": b,
      "stencil": b,
    }, ctx)

    canvas := NewCanvasRenderingContext2D(ctx)
    webgl := NewWebGLRenderingContext(ctx)

    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewLiteralString("2d", ctx), canvas},
      []values.Value{NewLiteralString("2d", ctx), o2d, canvas},
      []values.Value{NewLiteralString("webgl", ctx), webgl},
      []values.Value{NewLiteralString("webgl", ctx), ogl, webgl},
      []values.Value{s, canvas},
      []values.Value{s, o2d, canvas},
    }, ctx), nil
  case "height", "width":
    return f, nil
  case "toDataURL":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s}, // defaults: type="image/png", quality=0.92
      []values.Value{s, s},
      []values.Value{s, f, s},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *HTMLCanvasElement) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  f := NewNumber(ctx)

  switch key {
  case "height", "width":
    return f.Check(arg, ctx)
  default:
    return ctx.NewError("Error: HTMLCanvasElement." + key + " not setable")
  }
}

func (p *HTMLCanvasElement) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewHTMLCanvasElementPrototype(), ctx), nil
}
