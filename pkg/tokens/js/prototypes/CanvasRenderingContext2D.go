package prototypes

import (
  "../values"

  "../../context"
)

type CanvasRenderingContext2D struct {
  BuiltinPrototype
}

func NewCanvasRenderingContext2DPrototype() values.Prototype {
  return &CanvasRenderingContext2D{newBuiltinPrototype("CanvasRenderingContext2D")}
}

func NewCanvasRenderingContext2D(ctx context.Context) values.Value {
  return values.NewInstance(NewCanvasRenderingContext2DPrototype(), ctx)
}

func (p *CanvasRenderingContext2D) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*CanvasRenderingContext2D); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func (p *CanvasRenderingContext2D) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  f := NewNumber(ctx)
  s := NewString(ctx)

  switch key {
  case "direction", "fillStyle", "strokeStyle", "font", "globalCompositeOperator", "lineCap", "lineJoin", "shadowColor", "textAlign", "textBaseline", "globalAlpha", "lineWidth", "miterLimit", "shadowBlur", "shadowOffsetX", "shadowOffsetY":
    return nil, ctx.NewError("Error: is a setter only")
  case "arc":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f, f, f, f, nil},
      []values.Value{f, f, f, f, f, b, nil},
    }, ctx), nil
  case "arcTo":
    return values.NewFunction([]values.Value{f, f, f, f, f, nil}, ctx), nil
  case "beginPath", "closePath", "restore", "save", "stroke":
    return values.NewFunction([]values.Value{nil}, ctx), nil
  case "bezierCurveTo", "translate":
    return values.NewFunction([]values.Value{f, f, f, f, f, f, nil}, ctx), nil
  case "clearRect", "fillRect", "quadraticCurveTo", "rect", "strokeRect":
    return values.NewFunction([]values.Value{f, f, f, f, nil}, ctx), nil
  case "clip":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{nil},
      []values.Value{s, nil},
    }, ctx), nil
  case "createLinearGradient":
    return values.NewFunction([]values.Value{f, f, f, f, NewCanvasGradient(ctx)}, ctx), nil
  case "createRadialGradient":
    return values.NewFunction([]values.Value{f, f, f, f, f, f, NewCanvasGradient(ctx)}, ctx), nil
  case "createPattern":
    pat := NewCanvasPattern(ctx)

    return values.NewFunction([]values.Value{
      NewCanvasImageSource(ctx), s, pat,
    }, ctx), nil
  case "drawImage":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewCanvasImageSource(ctx), f, f, nil},
      []values.Value{NewCanvasImageSource(ctx), f, f, f, f, nil},
      []values.Value{NewCanvasImageSource(ctx), f, f, f, f, f, f, f, f, nil},
    }, ctx), nil
  case "ellipse":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f, f, f, f, f, f, nil},
      []values.Value{f, f, f, f, f, f, f, b, nil},
    }, ctx), nil
  case "fill":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{nil},
      []values.Value{s, nil},
    }, ctx), nil
  case "fillText", "strokeText":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{s, f, f, nil},
      []values.Value{s, f, f, f, nil},
    }, ctx), nil
  case "getImageData":
    return values.NewFunction([]values.Value{
      f, f, f, f, NewImageData(ctx),
    }, ctx), nil
  case "getTransform":
    return values.NewFunction([]values.Value{NewDOMMatrix(ctx)}, ctx), nil
  case "isPointInPath":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f, b}, 
      []values.Value{f, f, s, b}, 
    }, ctx), nil
  case "isPointInStroke":
    return values.NewFunction([]values.Value{f, f, b}, ctx), nil
  case "lineTo", "moveTo", "scale":
    return values.NewFunction([]values.Value{f, f, nil}, ctx), nil
  case "measureText":
    return values.NewFunction([]values.Value{s, NewTextMetrics(ctx)}, ctx), nil
  case "putImageData":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewImageData(ctx), f, f, nil},
      []values.Value{NewImageData(ctx), f, f, f, f, f, f, nil},
    }, ctx), nil
  case "rotate":
    return values.NewFunction([]values.Value{f, nil}, ctx), nil
  case "setTransform":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{f, f, f, f, f, f, nil},
      []values.Value{NewDOMMatrix(ctx), nil},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *CanvasRenderingContext2D) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  f := NewNumber(ctx)
  s := NewString(ctx)

  switch key {
  case "direction":
    return s.Check(arg, ctx)
  case "fillStyle", "strokeStyle":
    if !(IsString(arg) || IsCanvasPattern(arg) || IsCanvasGradient(arg)) {
      return ctx.NewError("Error: expected String, CanvasGradient or CanvasPattern, got " + arg.TypeName())
    } else {
      return nil
    }
  case "font", "globalCompositeOperator", "lineCap", "lineJoin", "shadowColor", "textAlign", "textBaseline":
    return s.Check(arg, ctx)
  case "globalAlpha", "lineWidth", "miterLimit", "shadowBlur", "shadowOffsetX", "shadowOffsetY":
    return f.Check(arg, ctx)
  default:
    return ctx.NewError("Error: " + p.Name() + "." + key + " is not a setable")
  }
}

func (p *CanvasRenderingContext2D) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewCanvasRenderingContext2DPrototype(), ctx), nil
}
